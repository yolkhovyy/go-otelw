package metricw

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_client "github.com/prometheus/client_model/go"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"go.opentelemetry.io/otel/metric"
)

// RegisterPrometheusCollectors maps metrics provided by Prometheus collectors to otel metrics.
// The metrics are sent via otlp http/grpc to Otel Collector.
func (m *Metric) RegisterPrometheusCollectors(ctx context.Context, colls ...prometheus.Collector) error {
	m.prometheusRegistry = prometheus.NewRegistry()
	m.prometheusRegistry.MustRegister(colls...)

	// Collect metrics once - to fetch their names, types,
	// and create corresponding instruments.
	err := m.gather(ctx, nil)
	if err != nil {
		return fmt.Errorf("register prometheus collectors: %w", err)
	}

	observables := make([]metric.Observable, 0,
		len(m.counters)+
			len(m.gauges))

	for _, counter := range m.counters {
		observables = append(observables, counter)
	}

	for _, gauge := range m.gauges {
		observables = append(observables, gauge)
	}

	// Register otel callback and observables.
	registration, err := m.meter.RegisterCallback(m.gather, observables...)
	if err != nil {
		return fmt.Errorf("register callback: %w", err)
	}

	m.registrations = append(m.registrations, registration)

	return nil
}

//nolint:gocognit,cyclop,funlen
func (m *Metric) gather(ctx context.Context, observer metric.Observer) error {
	metrics, err := m.prometheusRegistry.Gather()
	if err != nil {
		return fmt.Errorf("gather metrics: %w", err)
	}

	logger := slogw.DefaultLogger()

	for _, mFamily := range metrics {
		mfName := mFamily.GetName()
		mfType := mFamily.GetType()

		switch mfType {
		case prometheus_client.MetricType_GAUGE:
			instrument, exists := m.gauges[mfName]
			if !exists {
				var err error

				instrument, err = m.meter.Float64ObservableGauge(mfName, metric.WithDescription(mFamily.GetHelp()))
				if err != nil {
					return fmt.Errorf("create gauge instrument: %w", err)
				}

				m.gauges[mfName] = instrument
			}

			if observer != nil {
				for _, m := range mFamily.GetMetric() {
					mGauge := m.GetGauge()
					if mGauge != nil {
						mValue := mGauge.GetValue()
						observer.ObserveFloat64(instrument, mValue)
					}
				}
			}

		case prometheus_client.MetricType_COUNTER:
			instrument, exists := m.counters[mfName]
			if !exists {
				var err error

				instrument, err = m.meter.Float64ObservableCounter(mfName, metric.WithDescription(mFamily.GetHelp()))
				if err != nil {
					return fmt.Errorf("create gauge instrument: %w %s %s", ErrInvalidMetricType, mfName, mfType.String())
				}

				m.counters[mfName] = instrument
			}

			if observer != nil {
				for _, m := range mFamily.GetMetric() {
					mg := m.GetGauge()
					if mg != nil {
						mValue := mg.GetValue()
						observer.ObserveFloat64(instrument, mValue)
					}
				}
			}

		case prometheus_client.MetricType_HISTOGRAM:
			instrument, exists := m.histograms[mfName]
			if !exists {
				var err error

				instrument, err = m.meter.Float64Histogram(mfName, metric.WithDescription(mFamily.GetHelp()))
				if err != nil {
					return fmt.Errorf("create histogram instrument: %w", err)
				}

				m.histograms[mfName] = instrument
			}

			for _, m := range mFamily.GetMetric() {
				mHistogram := m.GetHistogram()
				if mHistogram != nil {
					instrument.Record(ctx, mHistogram.GetSampleSum())

					mBuckets := mHistogram.GetBucket()
					for _, bucket := range mBuckets {
						upperBound := bucket.GetUpperBound()
						count := bucket.GetCumulativeCountFloat()
						logger.DebugContext(ctx,
							"histogram bucket",
							slog.Float64("upper_bound", upperBound),
							slog.Float64("count", count),
						)
					}
				}
			}

		case prometheus_client.MetricType_SUMMARY:
			instrument, exists := m.histograms[mfName]
			if !exists {
				var err error

				instrument, err = m.meter.Float64Histogram(mfName, metric.WithDescription(mFamily.GetHelp()))
				if err != nil {
					return fmt.Errorf("create histogram instrument: %w", err)
				}

				m.histograms[mfName] = instrument
			}

			for _, m := range mFamily.GetMetric() {
				mSummary := m.GetSummary()
				if mSummary != nil {
					instrument.Record(ctx, mSummary.GetSampleSum())

					quantiles := mSummary.GetQuantile()
					for _, quantile := range quantiles {
						percentile := quantile.GetQuantile()
						value := quantile.GetValue()
						logger.DebugContext(ctx,
							"summary quantile",
							slog.Float64("quantile", percentile),
							slog.Float64("value", value),
						)
					}
				}
			}

		case prometheus_client.MetricType_GAUGE_HISTOGRAM, prometheus_client.MetricType_UNTYPED:
			logger.DebugContext(ctx,
				"unsupported prometheus metric type",
				slog.String("family name", mfName),
				slog.String("metric type", mfType.String()),
			)

		default:
			return fmt.Errorf("prometheus callback: %w %s %s", err, mfName, mfType.String())
		}
	}

	return nil
}

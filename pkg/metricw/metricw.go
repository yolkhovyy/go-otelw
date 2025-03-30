package metricw

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Metric represents a structured wrapper for OpenTelemetry metrics.
// It includes metric providers, exporters, and various metric types
// such as counters, gauges, and histograms.
type Metric struct {
	provider      *sdkmetric.MeterProvider
	exporter      sdkmetric.Exporter
	registrations []metric.Registration
	meter         metric.Meter

	counters   map[string]metric.Float64ObservableCounter
	gauges     map[string]metric.Float64ObservableGauge
	histograms map[string]metric.Float64Histogram

	prometheusRegistry *prometheus.Registry
}

// Configure initializes and configures the OpenTelemetry metric provider.
// It sets up the exporter, resource attributes, and meter provider.
// Optionally, it registers Prometheus collectors if enabled in the config.
func Configure( //nolint:cyclop,funlen
	ctx context.Context,
	config Config,
	attrs []attribute.KeyValue,
	writers ...io.Writer,
) (*Metric, error) {
	exporter, err := exporter(ctx, config, writers...)
	if err != nil {
		return nil, fmt.Errorf("metricw configure: %w", err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(attrs...))
	if err != nil {
		return nil, fmt.Errorf("metricw configure resource merge: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(config.Interval))

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(provider)

	serviceName := "undefined"

	attrSet := attribute.NewSet(attrs...)
	if value, exists := attrSet.Value(semconv.ServiceNameKey); exists {
		serviceName = value.AsString()
	}

	met := Metric{
		provider:      provider,
		exporter:      exporter,
		registrations: make([]metric.Registration, 0),
		gauges:        make(map[string]metric.Float64ObservableGauge),
		histograms:    make(map[string]metric.Float64Histogram),
		counters:      make(map[string]metric.Float64ObservableCounter),
		meter:         provider.Meter(serviceName),
	}

	if config.Prometheus {
		if err = met.RegisterPrometheusCollectors(ctx,
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		); err != nil {
			return nil, fmt.Errorf("metricw configure prometheus collectors: %w", err)
		}
	}

	return &met, nil
}

// Shutdown gracefully shuts down all metric-related components.
// It unregisters metric registrations, shuts down the provider,
// and ensures the exporter is properly closed.
func (m *Metric) Shutdown(ctx context.Context) error {
	var errs error

	for i := range m.registrations {
		err := m.registrations[i].Unregister()
		errs = errors.Join(errs, fmt.Errorf("metricw registration shutdown: %w", err))
	}

	if m.provider != nil {
		err := m.provider.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("metricw provider shutdown: %w", err))
	}

	if m.exporter != nil {
		err := m.exporter.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("metricw exporter shutdown: %w", err))
	}

	return errs
}

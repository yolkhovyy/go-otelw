package metricw

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/yolkhovyy/go-otelw/pkg/otelw"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

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

//nolint:cyclop,funlen
func Configure(
	ctx context.Context,
	config otelw.Metric,
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

	met := Metric{
		provider:      provider,
		exporter:      exporter,
		registrations: make([]metric.Registration, 0),
		gauges:        make(map[string]metric.Float64ObservableGauge),
		histograms:    make(map[string]metric.Float64Histogram),
		counters:      make(map[string]metric.Float64ObservableCounter),
		meter:         provider.Meter(otelw.FindAttribute(attrs, semconv.ServiceNameKey)),
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

//nolint:ireturn,cyclop,funlen
func exporter(
	ctx context.Context,
	config otelw.Metric,
	writers ...io.Writer,
) (sdkmetric.Exporter, error) {
	var err error

	var exporter sdkmetric.Exporter

	switch {
	case !config.Enable:
		exporter, err = stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
	case len(writers) > 0:
		options := make([]stdoutmetric.Option, 0, len(writers))
		for _, w := range writers {
			options = append(options, stdoutmetric.WithWriter(w))
		}

		exporter, err = stdoutmetric.New(options...)
	case config.Collector.Protocol == otelw.GRPC:
		options := []otlpmetricgrpc.Option{}
		if config.Collector.Connection != "" {
			options = append(options, otlpmetricgrpc.WithEndpoint(config.Collector.Connection))
		}

		if config.Collector.Insecure {
			options = append(options, otlpmetricgrpc.WithInsecure())
		} else {
			tlsConfig, err := otelw.TLSCredentials(config.Collector)
			if err != nil {
				return nil, fmt.Errorf("metricw otlp tls config: %w", err)
			}

			options = append(options, otlpmetricgrpc.WithTLSCredentials(tlsConfig))
		}

		exporter, err = otlpmetricgrpc.New(ctx, options...)
	case config.Collector.Protocol == otelw.HTTP:
		options := []otlpmetrichttp.Option{}
		if config.Collector.Connection != "" {
			options = append(options, otlpmetrichttp.WithEndpoint(config.Collector.Connection))
		}

		if config.Collector.Insecure {
			options = append(options, otlpmetrichttp.WithInsecure())
		} else {
			tlsConfig, err := otelw.TLSConfig(config.Collector)
			if err != nil {
				return nil, fmt.Errorf("metricw otlp tls config: %w", err)
			}

			options = append(options, otlpmetrichttp.WithTLSClientConfig(tlsConfig))
		}

		exporter, err = otlpmetrichttp.New(ctx, options...)
	default:
		err = fmt.Errorf("metricw exporter: %w %s", ErrInvalidProtocol, config.Collector.Protocol)
	}

	if err != nil {
		return nil, fmt.Errorf("metricw exporter: %w", err)
	}

	return exporter, err
}

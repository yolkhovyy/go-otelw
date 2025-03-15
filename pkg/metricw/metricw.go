package metricw

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/yolkhovyy/go-otelw/pkg/collector"
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

//nolint:ireturn
func exporter(
	ctx context.Context,
	config Config,
	writers ...io.Writer,
) (sdkmetric.Exporter, error) {
	var err error

	var exporter sdkmetric.Exporter

	switch {
	case !config.Enable:
		exporter, err = stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
	case len(writers) > 0:
		exporter, err = stdoutExporter(writers...)
	case config.Collector.Protocol == collector.GRPC:
		exporter, err = grpcExporter(ctx, config)
	case config.Collector.Protocol == collector.HTTP:
		exporter, err = httpExporter(ctx, config)
	default:
		err = fmt.Errorf("metricw exporter: %w %s", ErrInvalidProtocol, config.Collector.Protocol)
	}

	if err != nil {
		return nil, fmt.Errorf("metricw exporter: %w", err)
	}

	return exporter, err
}

//nolint:ireturn
func stdoutExporter(writers ...io.Writer) (sdkmetric.Exporter, error) {
	options := make([]stdoutmetric.Option, 0, len(writers))
	for _, w := range writers {
		options = append(options, stdoutmetric.WithWriter(w))
	}

	exporter, err := stdoutmetric.New(options...)
	if err != nil {
		return nil, fmt.Errorf("metricw stfout exporter: %w", err)
	}

	return exporter, nil
}

//nolint:ireturn
func grpcExporter(ctx context.Context, config Config) (sdkmetric.Exporter, error) {
	options := []otlpmetricgrpc.Option{}
	if config.Collector.Connection != "" {
		options = append(options, otlpmetricgrpc.WithEndpoint(config.Collector.Connection))
	}

	if config.Collector.Insecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	} else {
		tlsConfig, err := collector.TLSCredentials(config.Collector)
		if err != nil {
			return nil, fmt.Errorf("metricw otlp tls config: %w", err)
		}

		options = append(options, otlpmetricgrpc.WithTLSCredentials(tlsConfig))
	}

	exporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("metricw new exporter: %w", err)
	}

	return exporter, nil
}

//nolint:ireturn
func httpExporter(ctx context.Context, config Config) (sdkmetric.Exporter, error) {
	options := []otlpmetrichttp.Option{}
	if config.Collector.Connection != "" {
		options = append(options, otlpmetrichttp.WithEndpoint(config.Collector.Connection))
	}

	if config.Collector.Insecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	} else {
		tlsConfig, err := collector.TLSConfig(config.Collector)
		if err != nil {
			return nil, fmt.Errorf("metricw otlp tls config: %w", err)
		}

		options = append(options, otlpmetrichttp.WithTLSClientConfig(tlsConfig))
	}

	exporter, err := otlpmetrichttp.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("metricw new exporter: %w", err)
	}

	return exporter, nil
}

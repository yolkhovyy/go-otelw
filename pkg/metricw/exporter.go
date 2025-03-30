package metricw

import (
	"context"
	"fmt"
	"io"

	"github.com/yolkhovyy/go-otelw/pkg/collector"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// exporter selects and initializes a metric exporter based on the configuration.
// It supports different protocols (gRPC, HTTP, stdout) and returns an OpenTelemetry SDK exporter.
// If tracing is disabled, it returns a no-op exporter that discards all data.
func exporter( //nolint:ireturn
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

// stdoutExporter initializes an OpenTelemetry metric exporter that writes to the provided writers.
// It is used when metrics need to be logged to standard output or a file.
func stdoutExporter( //nolint:ireturn
	writers ...io.Writer,
) (sdkmetric.Exporter, error) {
	options := make([]stdoutmetric.Option, 0, len(writers))
	for _, w := range writers {
		options = append(options, stdoutmetric.WithWriter(w))
	}

	exporter, err := stdoutmetric.New(options...)
	if err != nil {
		return nil, fmt.Errorf("metricw stdout exporter: %w", err)
	}

	return exporter, nil
}

// grpcExporter initializes an OpenTelemetry metric exporter using gRPC protocol.
// It configures TLS credentials if necessary.
func grpcExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (sdkmetric.Exporter, error) {
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

// httpExporter initializes an OpenTelemetry metric exporter using HTTP protocol.
// It configures TLS settings if required.
func httpExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (sdkmetric.Exporter, error) {
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

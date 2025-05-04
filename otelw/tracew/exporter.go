package tracew

import (
	"context"
	"fmt"
	"io"

	"github.com/yolkhovyy/go-otelw/otelw/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// exporter initializes and returns a SpanExporter based on the provided configuration.
// It supports different protocols: stdout, gRPC, and HTTP, determined by config settings.
func exporter( //nolint:ireturn
	ctx context.Context,
	config Config,
	writers ...io.Writer,
) (sdktrace.SpanExporter, error) {
	var err error

	var exporter sdktrace.SpanExporter

	switch {
	case !config.Enable:
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	case len(writers) > 0:
		exporter, err = stdoutExporter(writers...)
	case config.OTLP.Protocol == otlp.GRPC:
		exporter, err = grpcExporter(ctx, config)
	case config.OTLP.Protocol == otlp.HTTP:
		exporter, err = httpExporter(ctx, config)
	default:
		err = fmt.Errorf("tracew exporter: %w %s", ErrInvalidProtocol, config.OTLP.Protocol)
	}

	return exporter, err
}

// stdoutExporter creates an exporter that writes traces to stdout or provided writers.
func stdoutExporter( //nolint:ireturn
	writers ...io.Writer,
) (sdktrace.SpanExporter, error) {
	options := make([]stdouttrace.Option, 0, len(writers))
	for _, w := range writers {
		options = append(options, stdouttrace.WithWriter(w))
	}

	exporter, err := stdouttrace.New(options...)
	if err != nil {
		return nil, fmt.Errorf("tracew stdout exporter: %w", err)
	}

	return exporter, nil
}

// grpcExporter initializes a gRPC-based OTLP exporter with optional TLS settings.
func grpcExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (sdktrace.SpanExporter, error) {
	options := []otlptracegrpc.Option{}
	if config.OTLP.Endpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.OTLP.Endpoint))
	}

	if config.OTLP.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	} else {
		tslCreds, err := otlp.TLSCredentials(config.OTLP)
		if err != nil {
			return nil, fmt.Errorf("tracew otlp grpc tls credentials: %w", err)
		}

		options = append(options, otlptracegrpc.WithTLSCredentials(tslCreds))
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("tracew new exporter: %w", err)
	}

	return exporter, nil
}

// httpExporter initializes an HTTP-based OTLP exporter with optional TLS settings.
func httpExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (sdktrace.SpanExporter, error) {
	options := []otlptracehttp.Option{}
	if config.OTLP.Endpoint != "" {
		options = append(options, otlptracehttp.WithEndpoint(config.OTLP.Endpoint))
	}

	if config.OTLP.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	} else {
		tlsConfig, err := otlp.TLSConfig(config.OTLP)
		if err != nil {
			return nil, fmt.Errorf("tracew otlp http tls config: %w", err)
		}

		options = append(options, otlptracehttp.WithTLSClientConfig(tlsConfig))
	}

	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("tracew new exporter: %w", err)
	}

	return exporter, nil
}

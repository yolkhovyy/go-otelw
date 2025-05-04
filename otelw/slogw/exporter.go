package slogw

import (
	"context"
	"fmt"
	"io"

	"github.com/yolkhovyy/go-otelw/otelw/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
)

// exporter initializes and returns a log.Exporter based on the given configuration and writers.
// If the exporter setup fails, it returns an error.
func exporter( //nolint:ireturn
	ctx context.Context,
	config Config,
	writers ...io.Writer,
) (log.Exporter, error) {
	var err error

	var exporter log.Exporter

	switch {
	case !config.Enable:
		exporter, err = stdoutlog.New(stdoutlog.WithWriter(io.Discard))
	case len(writers) > 0:
		exporter, err = stdoutExporter(writers...)
	case config.OTLP.Protocol == otlp.GRPC:
		exporter, err = grpcExporter(ctx, config)
	case config.OTLP.Protocol == otlp.HTTP:
		exporter, err = httpExporter(ctx, config)
	default:
		err = fmt.Errorf("slogw exporter: %w %s", ErrInvalidProtocol, config.OTLP.Protocol)
	}

	if err != nil {
		return nil, fmt.Errorf("slogw exporter: %w", err)
	}

	return exporter, nil
}

// stdoutExporter initializes and returns a stdout-based log.Exporter for the given writers.
// If the exporter setup fails, it returns an error.
func stdoutExporter( //nolint:ireturn
	writers ...io.Writer,
) (log.Exporter, error) {
	options := make([]stdoutlog.Option, 0, len(writers))
	for _, w := range writers {
		options = append(options, stdoutlog.WithWriter(w))
	}

	exporter, err := stdoutlog.New(options...)
	if err != nil {
		return nil, fmt.Errorf("slogw stdout exporter: %w", err)
	}

	return exporter, nil
}

// grpcExporter initializes and returns a gRPC-based log.Exporter using the provided configuration.
// If the exporter setup fails, it returns an error.
func grpcExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (log.Exporter, error) {
	options := []otlploggrpc.Option{}
	if config.OTLP.Endpoint != "" {
		options = append(options, otlploggrpc.WithEndpoint(config.OTLP.Endpoint))
	}

	if config.OTLP.Insecure {
		options = append(options, otlploggrpc.WithInsecure())
	} else {
		tslCreds, err := otlp.TLSCredentials(config.OTLP)
		if err != nil {
			return nil, fmt.Errorf("slogw otlp tls credentials: %w", err)
		}

		options = append(options, otlploggrpc.WithTLSCredentials(tslCreds))
	}

	exporter, err := otlploggrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("slogw new exporter: %w", err)
	}

	return exporter, nil
}

// httpExporter initializes and returns an HTTP-based log.Exporter using the provided configuration.
// If the exporter setup fails, it returns an error.
func httpExporter( //nolint:ireturn
	ctx context.Context,
	config Config,
) (log.Exporter, error) {
	options := []otlploghttp.Option{}
	if config.OTLP.Endpoint != "" {
		options = append(options, otlploghttp.WithEndpoint(config.OTLP.Endpoint))
	}

	if config.OTLP.Insecure {
		options = append(options, otlploghttp.WithInsecure())
	} else {
		tlsConfig, err := otlp.TLSConfig(config.OTLP)
		if err != nil {
			return nil, fmt.Errorf("slogw otlp tls config: %w", err)
		}

		options = append(options, otlploghttp.WithTLSClientConfig(tlsConfig))
	}

	exporter, err := otlploghttp.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("slogw new exporter: %w", err)
	}

	return exporter, nil
}

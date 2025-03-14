package slogw

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/yolkhovyy/go-otelw/pkg/collector"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Logger struct {
	*slog.Logger
	exporter log.Exporter
	provider *log.LoggerProvider
}

//nolint:cyclop,funlen
func Configure(
	ctx context.Context,
	config Config,
	attrs []attribute.KeyValue,
	writers ...io.Writer,
) (*Logger, error) {
	exporter, err := exporter(ctx, config, writers...)
	if err != nil {
		return nil, fmt.Errorf("slogw configure: %w", err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(attrs...))
	if err != nil {
		return nil, fmt.Errorf("slogw configure resource merge: %w", err)
	}

	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(&WithSeverityText{exporter})),
	)

	serviceName := "undefined"

	attrSet := attribute.NewSet(attrs...)
	if value, exists := attrSet.Value(semconv.ServiceNameKey); exists {
		serviceName = value.AsString()
	}

	// TODO: check other options.
	slog.SetDefault(otelslog.NewLogger(
		serviceName,
		otelslog.WithLoggerProvider(provider),
		otelslog.WithSource(config.Caller),
	))

	return &Logger{
		Logger:   slog.Default(),
		exporter: exporter,
		provider: provider,
	}, nil
}

func (l *Logger) Shutdown(ctx context.Context) error {
	var errs error

	if l.provider != nil {
		err := l.provider.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("slogw provider shutdown: %w", err))
	}

	if l.exporter != nil {
		err := l.exporter.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("slogw exporter shutdown: %w", err))
	}

	return errs
}

func (l *Logger) ForceFlush(ctx context.Context) error {
	if err := l.provider.ForceFlush(ctx); err != nil {
		return fmt.Errorf("slogw force flush: %w", err)
	}

	return nil
}

func NewLogger() *slog.Logger {
	return slog.New(slog.Default().Handler())
}

//nolint:ireturn,cyclop,funlen
func exporter(
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
		options := make([]stdoutlog.Option, 0, len(writers))
		for _, w := range writers {
			options = append(options, stdoutlog.WithWriter(w))
		}

		exporter, err = stdoutlog.New(options...)
	case config.Collector.Protocol == collector.GRPC:
		options := []otlploggrpc.Option{}
		if config.Collector.Connection != "" {
			options = append(options, otlploggrpc.WithEndpoint(config.Collector.Connection))
		}

		if config.Collector.Insecure {
			options = append(options, otlploggrpc.WithInsecure())
		} else {
			tslCreds, err := collector.TLSCredentials(config.Collector)
			if err != nil {
				return nil, fmt.Errorf("slogw otlp tls credentials: %w", err)
			}

			options = append(options, otlploggrpc.WithTLSCredentials(tslCreds))
		}

		exporter, err = otlploggrpc.New(ctx, options...)
	case config.Collector.Protocol == collector.HTTP:
		options := []otlploghttp.Option{}
		if config.Collector.Connection != "" {
			options = append(options, otlploghttp.WithEndpoint(config.Collector.Connection))
		}

		if config.Collector.Insecure {
			options = append(options, otlploghttp.WithInsecure())
		} else {
			tlsConfig, err := collector.TLSConfig(config.Collector)
			if err != nil {
				return nil, fmt.Errorf("slogw otlp tls config: %w", err)
			}

			options = append(options, otlploghttp.WithTLSClientConfig(tlsConfig))
		}

		exporter, err = otlploghttp.New(ctx, options...)
	default:
		err = fmt.Errorf("slogw exporter: %w %s", ErrInvalidProtocol, config.Collector.Protocol)
	}

	if err != nil {
		return nil, fmt.Errorf("slogw exporter: %w", err)
	}

	return exporter, nil
}

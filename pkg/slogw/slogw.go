package slogw

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Logger is a wrapper around the standard slog.Logger, providing additional
// functionality for exporting logs and managing a logger provider.
// It integrates with OpenTelemetry to support structured logging with trace
// and span information.
//
// Fields:
// - Logger: The embedded slog.Logger instance used for logging operations.
// - exporter: The log.Exporter responsible for exporting log data to the configured destination.
// - provider: The log.LoggerProvider that manages the lifecycle and configuration of the logger.
type Logger struct {
	*slog.Logger
	exporter log.Exporter
	provider *log.LoggerProvider
}

// Configure sets up a new default global Logger with the given configuration, attributes, and writers.
// It returns the configured Logger or an error if the setup fails.
func Configure( //nolint:cyclop,funlen
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

// Shutdown gracefully shuts down the Logger, ensuring all logs are flushed.
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

// ForceFlush forces the Logger to flush all buffered logs.
func (l *Logger) ForceFlush(ctx context.Context) error {
	if err := l.provider.ForceFlush(ctx); err != nil {
		return fmt.Errorf("slogw force flush: %w", err)
	}

	return nil
}

// NewLogger creates and returns a new instance of slog.Logger.
func NewLogger() *slog.Logger {
	return slog.New(slog.Default().Handler())
}

// DefaultLogger returns a default instance of slog.Logger.
func DefaultLogger() *slog.Logger {
	return slog.Default()
}

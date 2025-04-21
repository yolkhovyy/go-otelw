package tracew

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Tracer is a wrapper around the OpenTelemetry TracerProvider and SpanExporter.
// It manages the lifecycle of tracing components and provides methods for configuration and shutdown.
type Tracer struct {
	provider *sdktrace.TracerProvider
	exporter sdktrace.SpanExporter
}

// Configure sets up the Tracer with the given configuration, attributes, and optional writers.
// It initializes the TracerProvider and SpanExporter, and configures the global OpenTelemetry settings.
func Configure(
	ctx context.Context,
	config Config,
	attrs []attribute.KeyValue,
	writers ...io.Writer,
) (*Tracer, error) {
	exporter, err := exporter(ctx, config, writers...)
	if err != nil {
		return nil, fmt.Errorf("tracew configure: %w", err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(attrs...))
	if err != nil {
		return nil, fmt.Errorf("tracew configure resource merge: %w", err)
	}

	extras, err := resource.New(
		ctx,
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("tracew configure extra resources merge: %w", err)
	}

	res, err = resource.Merge(res, extras)
	if err != nil {
		return nil, fmt.Errorf("tracew configure resource merge: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)))

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(provider)

	return &Tracer{
		provider: provider,
		exporter: exporter,
	}, nil
}

// Shutdown gracefully shuts down the Tracer, ensuring all spans are flushed and resources are released.
// It returns any errors encountered during the shutdown process.
func (t *Tracer) Shutdown(ctx context.Context) error {
	var errs error

	if t.provider != nil {
		err := t.provider.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("tracew provider shutdown: %w", err))
	}

	if t.exporter != nil {
		err := t.exporter.Shutdown(ctx)
		errs = errors.Join(errs, fmt.Errorf("tracew exporter shutdown: %w", err))
	}

	return errs
}

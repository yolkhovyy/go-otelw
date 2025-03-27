package tracew

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/yolkhovyy/go-otelw/pkg/collector"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Tracer struct {
	provider *sdktrace.TracerProvider
	exporter sdktrace.SpanExporter
}

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

//nolint:ireturn,cyclop
func exporter(
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
	case config.Collector.Protocol == collector.GRPC:
		exporter, err = grpcExporter(ctx, config)
	case config.Collector.Protocol == collector.HTTP:
		exporter, err = httpExporter(ctx, config)
	default:
		err = fmt.Errorf("tracew exporter: %w %s", ErrInvalidProtocol, config.Collector.Protocol)
	}

	return exporter, err
}

//nolint:ireturn
func stdoutExporter(writers ...io.Writer) (sdktrace.SpanExporter, error) {
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

//nolint:ireturn
func grpcExporter(ctx context.Context, config Config) (sdktrace.SpanExporter, error) {
	options := []otlptracegrpc.Option{}
	if config.Collector.Connection != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.Collector.Connection))
	}

	if config.Collector.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	} else {
		tslCreds, err := collector.TLSCredentials(config.Collector)
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

//nolint:ireturn
func httpExporter(ctx context.Context, config Config) (sdktrace.SpanExporter, error) {
	options := []otlptracehttp.Option{}
	if config.Collector.Connection != "" {
		options = append(options, otlptracehttp.WithEndpoint(config.Collector.Connection))
	}

	if config.Collector.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	} else {
		tlsConfig, err := collector.TLSConfig(config.Collector)
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

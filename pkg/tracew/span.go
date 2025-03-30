package tracew

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Span wraps trace.Span to provide additional functionality.
type Span struct {
	trace.Span
}

// Start initializes and starts a new span using OpenTelemetry tracing.
// It returns a new context containing the span and a wrapped Span struct.
// The function accepts a tracer name (tname), span name (sname),
// and optional span start options.
func Start(
	ctx context.Context,
	tname, sname string,
	options ...trace.SpanStartOption,
) (context.Context, Span) {
	ctx, span := otel.Tracer(tname).Start(ctx, sname, options...) //nolint:spancheck

	return ctx, Span{Span: span} //nolint:spancheck
}

// End finalizes the span. If an error is provided, it records the error
// and marks the span status as an error; otherwise, it marks it as successful.
// It accepts optional span end options.
func (s *Span) End(err error, options ...trace.SpanEndOption) {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
	} else {
		s.SetStatus(codes.Ok, "No error")
	}

	s.Span.End(options...)
}

package tracew

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	trace.Span
}

//nolint:spancheck
func Start(
	ctx context.Context,
	tname, sname string,
	options ...trace.SpanStartOption,
) (context.Context, Span) {
	ctx, span := otel.Tracer(tname).Start(ctx, sname, options...)

	return ctx, Span{Span: span}
}

func (s *Span) End(err error, options ...trace.SpanEndOption) {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
	} else {
		s.SetStatus(codes.Ok, "No error")
	}

	s.Span.End(options...)
}

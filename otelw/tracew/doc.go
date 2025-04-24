// provides a convenience wrapper around the
// OpenTelemetry SDK trace (see go.opentelemetry.io/otel/sdk/trace) package.
//
// simplifies tracing configuration and shutdown.
//
// Usage:
//
//	tracer, err := tracew.Configure(ctx, config.Logger, attrs, writers...)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "tracer configure: %v", err)
//		return osx.ExitFailure
//	}
//	defer func() {
//		err := errors.Join(err, tracer.Shutdown(ctx))
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "tracer shutdown: %v", err)
//		}
//	}()
//
//	ctx, span := tracew.Start(ctx, "echo", "worker",
//		trace.WithAttributes(attribute.Int("sequence", sequence)),
//	)
//	defer func() { span.End(err) }()
package tracew

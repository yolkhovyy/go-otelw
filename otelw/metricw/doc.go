// Package tracew provides a convenience wrapper around the
// OpenTelemetry SDK `metric` (see go.opentelemetry.io/otel/sdk/metric) package.
//
// This wrapper simplifies metric configuration and shutdown.
//
// Usage:
//
//	metric, err := metricw.Configure(ctx, config.Logger, attrs, writers...)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "metric configure: %v", err)
//		return osx.ExitFailure
//	}
//	defer func() {
//		err := errors.Join(err, metric.Shutdown(ctx))
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "metric shutdown: %v", err)
//		}
//	}()
package metricw

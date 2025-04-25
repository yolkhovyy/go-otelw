// Provides a convenience wrapper around the
// log/slog (see https://pkg.go.dev/log/slog) package.
//
// Simplifies logging configuration and shutdown, and provides
// interoperability with OpenTelemetry.
//
// Usage:
//
//	logger, err := slogw.Configure(ctx, config.Logger, attrs, writers...)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "log configure: %v", err)
//		return osx.ExitFailure
//	}
//	defer func() {
//		err := errors.Join(err, logger.Shutdown(ctx))
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "slog shutdown: %v", err)
//		}
//	}()
//
//	logger.InfoContext(ctx, "build info",
//		slog.String("version", version.Tag),
//		slog.String("time", buildInfo.Time),
//		slog.String("revision", buildInfo.Revision),
package slogw

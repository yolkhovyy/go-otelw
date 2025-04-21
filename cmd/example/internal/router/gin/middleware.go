package gin

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
)

// Logger middleware logs HTTP requests.
func Logger() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		start := time.Now()

		logger := slogw.DefaultLogger()

		logger.InfoContext(
			gctx.Request.Context(),
			"http request begin",
			slog.String("method", gctx.Request.Method),
			slog.String("path", gctx.Request.URL.Path),
		)

		gctx.Next()

		logger.InfoContext(
			gctx.Request.Context(),
			"http request end",
			slog.String("method", gctx.Request.Method),
			slog.String("path", gctx.Request.URL.Path),
			slog.Int("status", gctx.Writer.Status()),
			slog.Duration("duration(ms)", time.Since(start).Round(time.Millisecond)),
		)
	}
}

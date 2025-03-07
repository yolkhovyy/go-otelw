package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/yolkhovyy/go-otelw/cmd/example/internal/domain"
	"github.com/yolkhovyy/go-otelw/cmd/example/internal/otelw"
	ginrouter "github.com/yolkhovyy/go-otelw/cmd/example/internal/router/gin"
	httpserver "github.com/yolkhovyy/go-otelw/cmd/example/internal/server/http"
	"github.com/yolkhovyy/go-utilities/buildinfo"
	"github.com/yolkhovyy/go-utilities/osx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	serviceName  = "example"
	configPrefix = "EXAMPLE"
)

func main() {
	os.Exit(run())
}

//nolint:funlen
func run() osx.ExitCode {
	// Build info.
	buildInfo := buildinfo.ReadData()

	// Configuration.
	configFile := flag.String("config", "config.yml", "Path to the configuration file (default: config.yml)")
	flag.Parse()

	config := Config{}
	if err := config.Load(*configFile, configPrefix); err != nil {
		fmt.Fprintf(os.Stderr, "config load: %v", err)

		return osx.ExitConfigError
	}

	// Context.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Telemetry.
	serviceAttributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(buildInfo.Version),
	}

	logger, tracer, metric, err := otelw.Configure(ctx, config.Config, serviceAttributes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otelw configure: %v", err)

		return osx.ExitFailure
	}

	defer func() {
		err := errors.Join(err,
			metric.Shutdown(ctx),
			tracer.Shutdown(ctx),
			logger.Shutdown(ctx))
		if err != nil {
			fmt.Fprintf(os.Stderr, "otelw shutdown: %v", err)
		}
	}()

	logger.InfoContext(ctx, "build info",
		slog.String("version", buildInfo.Version),
		slog.String("time", buildInfo.Time),
		slog.String("commit", buildInfo.Commit),
	)

	// Domain.
	domain, err := domain.New(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "domain",
			slog.String("new", err.Error()),
		)

		return osx.ExitFailure
	}

	defer func() {
		if err := domain.Close(); err != nil {
			logger.ErrorContext(ctx, "domain",
				slog.String("close", err.Error()),
			)
		}
	}()

	// Router.
	router := ginrouter.New(domain,
		gin.Recovery(), ginrouter.Logger(), otelgin.Middleware(serviceName))

	// HTTP server.
	server := httpserver.New(config.HTTP, router.Handler())
	if err := server.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "http server",
			slog.String("run", err.Error()),
		)

		return osx.ExitFailure
	}

	logger.InfoContext(ctx, "shutting down")

	return osx.ExitSuccess
}

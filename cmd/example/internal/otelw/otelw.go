package otelw

import (
	"context"
	"fmt"
	"io"

	"github.com/yolkhovyy/go-otelw/pkg/metricw"
	"github.com/yolkhovyy/go-otelw/pkg/otelw"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
	"go.opentelemetry.io/otel/attribute"
)

func Configure(
	ctx context.Context,
	config otelw.Config,
	attrs []attribute.KeyValue,
	writers ...io.Writer,
) (*slogw.Logger, *tracew.Tracer, *metricw.Metric, error) {
	logger, err := slogw.Configure(ctx, config.Logger, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("slogw configure: %w", err)
	}

	tracer, err := tracew.Configure(ctx, config.Tracer, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("tracew configure: %w", err)
	}

	metric, err := metricw.Configure(ctx, config.Metric, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("metricw configure: %w", err)
	}

	return logger, tracer, metric, nil
}

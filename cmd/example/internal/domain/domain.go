package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrTimeout = errors.New("timeout")

func (u Controller) Echo(ctx context.Context, input string, count int) (string, error) {
	outChan := make(chan string, count)
	errChan := make(chan error, count)

	waitGroup := sync.WaitGroup{}

	for i := range count {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			worker(ctx, i, input, outChan, errChan)
		}()
	}

	go func() {
		waitGroup.Wait()
		close(outChan)
		close(errChan)
	}()

	var errs error

	outs := make([]string, 0, count)

	for {
		var err error

		var out string

		var outOk, errOk bool

		select {
		case out, outOk = <-outChan:
			if outOk {
				outs = append(outs, out)
			}

		case err, errOk = <-errChan:
			if errOk {
				errs = errors.Join(err, errs)
			}
		}

		if !outOk && !errOk {
			break
		}
	}

	return strings.Join(outs, ", ") + "...\n", errs
}

func worker(
	ctx context.Context,
	sequence int,
	input string,
	outChan chan<- string,
	errChan chan<- error,
) {
	var err error

	ctx, span := tracew.Start(ctx, "echo", "worker",
		trace.WithAttributes(attribute.Int("sequence", sequence)),
	)
	defer func() { span.End(err) }()

	logger := slogw.DefaultLogger()

	const workThreshold = 10

	// Do worker's work.
	{
		time.Sleep(time.Duration(sequence+1) * time.Millisecond)
	}

	msg := "do echo"

	logAttrs := []any{
		slog.Int("sequence", sequence),
		slog.String("input", input),
	}

	eventAttrs := []attribute.KeyValue{
		attribute.Int("sequence", sequence),
		attribute.String("input", input),
	}

	if sequence > workThreshold {
		err = fmt.Errorf("%s: %w", msg, ErrTimeout)
		logger.ErrorContext(ctx, msg, append(logAttrs, slog.String("error", err.Error()))...)
		span.AddEvent(msg, trace.WithAttributes(append(eventAttrs, attribute.String("error", err.Error()))...))
		errChan <- err
	} else {
		logger.InfoContext(ctx, msg, logAttrs...)
		span.AddEvent(msg, trace.WithAttributes(eventAttrs...))
		outChan <- input
	}
}

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
)

var ErrTimeout = errors.New("timeout")

//nolint:cyclop
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

	outputs := make([]string, 0, count)

	for {
		var outOpen, errOpen bool

		select {
		case output, ok := <-outChan:
			outOpen = ok
			if outOpen {
				outputs = append(outputs, output)
			}

		case err, ok := <-errChan:
			errOpen = ok

			switch {
			case !errOpen:
			case errs == nil:
				errs = err
			default:
				errs = fmt.Errorf("%w: %w", errs, err)
			}
		}

		if !outOpen && !errOpen {
			break
		}
	}

	return strings.Join(outputs, ", ") + "...\n", errs
}

func worker(
	ctx context.Context,
	sequence int,
	input string,
	outChan chan<- string,
	errChan chan<- error,
) {
	var err error

	ctx, span := tracew.Start(ctx, "echo", "worker")
	defer func() { span.End(err) }()

	logger := slogw.NewLogger()

	// Do work.
	const workThreshold = 10

	time.Sleep(time.Duration(sequence+1) * time.Millisecond)

	if sequence > workThreshold {
		err = fmt.Errorf("worker: %w", ErrTimeout)
		span.AddEvent(fmt.Errorf("sequence: %d input: %s error:%w", sequence, input, err).Error())
		logger.InfoContext(ctx, "echo worker",
			slog.Int("sequence", sequence),
			slog.String("input", input),
		)
		errChan <- err
	} else {
		span.AddEvent(fmt.Sprintf("sequence: %d input: %s", sequence, input))
		logger.InfoContext(ctx, "echo worker",
			slog.Int("sequence", sequence),
			slog.String("input", input),
		)
		outChan <- input
	}
}

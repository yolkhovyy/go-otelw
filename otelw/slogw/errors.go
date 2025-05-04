package slogw

import "errors"

var (
	// ErrInvalidFormat is returned when config.Format is not qual to slogw.Console or slogw.JSON.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrInvalidProtocol is returned when config.Collector.Protocol is not qual to otlp.GRPC or otlp.HTTP.
	ErrInvalidProtocol = errors.New("invalid protocol")
)

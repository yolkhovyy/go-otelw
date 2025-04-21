package metricw

import "errors"

var (
	// ErrInvalidFormat is returned when config.Format is not qual to slogw.Console or slogw.JSON.
	ErrInvalidProtocol = errors.New("invalid protocol")

	// ErrInvalidMetricType is returned when a not supported Prometheus metric type is requested.
	ErrInvalidMetricType = errors.New("invalid metric type")
)

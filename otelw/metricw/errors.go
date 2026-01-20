package metricw

import "errors"

var (
	// ErrInvalidProtocol is returned when config.Protocol is not equal to otlp.GRPC or otlp.HTTP.
	ErrInvalidProtocol = errors.New("invalid protocol")

	// ErrInvalidMetricType is returned when a not supported Prometheus metric type is requested.
	ErrInvalidMetricType = errors.New("invalid metric type")
)

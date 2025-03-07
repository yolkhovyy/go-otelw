package metricw

import "errors"

var (
	ErrInvalidProtocol   = errors.New("invalid protocol")
	ErrInvalidMetricType = errors.New("invalid metric type")
)

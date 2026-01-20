package tracew

import "errors"

// ErrInvalidProtocol is returned when config.Protocol is not equal to otlp.GRPC or otlp.HTTP.
var ErrInvalidProtocol = errors.New("invalid protocol")

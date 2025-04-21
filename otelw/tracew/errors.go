package tracew

import "errors"

// ErrInvalidFormat is returned when config.Format is not qual to slogw.Console or slogw.JSON.
var ErrInvalidProtocol = errors.New("invalid protocol")

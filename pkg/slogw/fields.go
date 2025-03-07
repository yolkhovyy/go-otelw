package slogw

import (
	"errors"
	"fmt"

	"github.com/yolkhovyy/go-utilities/stringx"
)

type FieldName string

const (
	FieldNameSpanID  FieldName = "spanId"
	FieldNameTraceID FieldName = "traceId"
)

type Format string

const (
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrTypeAssertion = errors.New("type assertion")
	ErrInvalidType   = errors.New("invalid type")
)

func (f *Format) String() string {
	return string(*f)
}

func (f *Format) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strFormat string
	if err := unmarshal(&strFormat); err != nil {
		return err
	}

	format := Format(stringx.TrimSpaceToLower(strFormat))
	switch format {
	case FormatConsole, FormatJSON:
		*f = format
	default:
		return fmt.Errorf("unmarshal: %w %s", ErrInvalidFormat, strFormat)
	}

	return nil
}

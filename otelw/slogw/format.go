package slogw

import (
	"fmt"

	"github.com/yolkhovyy/go-utilities/stringx"
)

// Format of logging.
type Format string

const (
	Console Format = "console"
	JSON    Format = "json"
)

// String returns the string representation of the Format.
func (f *Format) String() string {
	return string(*f)
}

// UnmarshalYAML unmarshals a YAML value into a Format, validating the format.
// It returns an error if the format is invalid.
func (f *Format) UnmarshalYAML(unmarshal func(any) error) error {
	var strFormat string
	if err := unmarshal(&strFormat); err != nil {
		return err
	}

	format := Format(stringx.TrimSpaceToLower(strFormat))
	switch format {
	case Console, JSON:
		*f = format
	default:
		return fmt.Errorf("unmarshal: %w %s", ErrInvalidFormat, strFormat)
	}

	return nil
}

package slogw

import (
	"time"

	"github.com/yolkhovyy/go-otelw/otelw/otlp"
)

// Config holds the configuration settings for the slogw package.
// It includes options for enabling logging, specifying the log format,
// setting the log level, and configuring the otlp protocol.
type Config struct {
	// Enable indicates whether logging is enabled.
	Enable bool `json:"enable" yaml:"enable" mapstructure:"enable"`

	// Caller specifies whether to include caller information in logs.
	Caller bool `json:"caller" yaml:"caller" mapstructure:"Caller"`

	// Format defines the output format of the logs - json (default), console.
	Format Format `json:"format" yaml:"format" mapstructure:"Format"`

	// Level sets the minimum log level - error, warn, info (default), debug.
	Level string `json:"level" yaml:"level" mapstructure:"Level"`

	// TimeFormat specifies the format for timestamps in logs.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat" mapstructure:"TimeFormat"`

	// OTLP holds the configuration for the OTEL protocol.
	OTLP otlp.Config `json:"otlp" yaml:"otlp" mapstructure:"OTLP"`
}

// Defaults returns a map of default configuration values for the slogw package.
// It includes default settings for enabling logging, caller information,
// log format, log level, and time format. It also incorporates defaults
// for the otlp package.
func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable
	defaults["Caller"] = DefaultCaller
	defaults["Format"] = DefaultFormat
	defaults["Level"] = DefaultLevel
	defaults["TimeFormat"] = DefaultTimeFormat

	for k, v := range otlp.Defaults() {
		defaults["OTLP."+k] = v
	}

	return defaults
}

const (
	// DefaultEnable is the default setting for enabling logging.
	DefaultEnable = false

	// DefaultCaller is the default setting for including caller information in logs.
	DefaultCaller = false

	// DefaultFormat is the default log output format.
	DefaultFormat = JSON

	// DefaultLevel is the default minimum log level.
	DefaultLevel = "info"

	// DefaultTimeFormat is the default format for timestamps in logs.
	DefaultTimeFormat = time.RFC3339
)

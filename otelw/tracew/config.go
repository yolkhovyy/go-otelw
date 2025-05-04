package tracew

import "github.com/yolkhovyy/go-otelw/otelw/otlp"

// Config holds the configuration settings for the tracew package.
type Config struct {
	// Enable indicates whether tracings is enabled.
	Enable bool `json:"enable" yaml:"enable" mapstructure:"enable"`

	// OTLP holds the configuration for the OTEL protocol.
	OTLP otlp.Config `json:"otlp" yaml:"otlp" mapstructure:"otlp"`
}

// Defaults returns a map of default configuration values for the tracew package.
// It includes default settings for enabling tracing and defaults
// for the otlp package.
func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable

	for k, v := range otlp.Defaults() {
		defaults["OTLP."+k] = v
	}

	return defaults
}

// DefaultEnable defines whether tracing is enabled by default.
const DefaultEnable = false

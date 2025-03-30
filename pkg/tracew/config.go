package tracew

import "github.com/yolkhovyy/go-otelw/pkg/collector"

// Config holds the configuration settings for the tracew package.
type Config struct {
	// Enable indicates whether tracings is enabled.
	Enable bool `yaml:"enable" mapstructure:"enable"`

	// Collector holds the configuration for the OTEL collector.
	Collector collector.Config `yaml:"collector" mapstructure:"collector"`
}

// Defaults returns a map of default configuration values for the tracew package.
// It includes default settings for enabling tracing and defaults
// for the collector package.
func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable

	for k, v := range collector.Defaults() {
		defaults["Collector."+k] = v
	}

	return defaults
}

// DefaultEnable defines whether tracing is enabled by default.
const DefaultEnable = false

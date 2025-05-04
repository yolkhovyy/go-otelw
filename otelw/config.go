package otelw

import (
	"github.com/yolkhovyy/go-otelw/otelw/metricw"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"github.com/yolkhovyy/go-otelw/otelw/tracew"
)

// Config defines the configuration structure for the OpenTelemetry wrapper.
// It includes configurations for logging, tracing, and metrics.
type Config struct {
	// Logging configuration
	Logger slogw.Config `json:"logger" yaml:"logger" mapstructure:"Logger"`

	// Tracing configuration
	Tracer tracew.Config `json:"tracer" yaml:"tracer" mapstructure:"Tracer"`

	// Metrics configuration
	Metric metricw.Config `json:"metric" yaml:"metric" mapstructure:"Metric"`
}

// Defaults returns a map containing default configuration values
// for the logger, tracer, and metrics components.
func Defaults() map[string]any {
	defaults := make(map[string]any)

	// Load default logger configuration
	for k, v := range slogw.Defaults() {
		defaults["Logger."+k] = v
	}

	// Load default tracer configuration
	for k, v := range tracew.Defaults() {
		defaults["Tracer."+k] = v
	}

	// Load default metrics configuration
	for k, v := range metricw.Defaults() {
		defaults["Metric."+k] = v
	}

	return defaults
}

package metricw

import (
	"time"

	"github.com/yolkhovyy/go-otelw/otelw/collector"
)

// Config holds the configuration settings for the package metricw package.
type Config struct {
	// Enable indicates whether metrics are enabled.
	Enable bool `yaml:"enable" mapstructure:"enable"`

	// Prometheus indicates whether Prometheus metric mapping is enabled.
	Prometheus bool `yaml:"prometheus" mapstructure:"prometheus"`

	// Interval holds metric collection interval.
	Interval time.Duration `yaml:"interval" mapstructure:"interval"`

	// Collector holds the configuration for the OTEL collector.
	Collector collector.Config `yaml:"collector" mapstructure:"collector"`
}

// Defaults returns a map of default configuration values for the metricw package.
// It includes default settings for enabling metrics, prometheus metrics mapping,
// metrics collection interval and defaults for the collector package.
func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable
	defaults["Prometheus"] = DefaultPrometheus
	defaults["Interval"] = DefaultInterval

	for k, v := range collector.Defaults() {
		defaults["Collector."+k] = v
	}

	return defaults
}

const (
	// DefaultEnable is the default setting for enabling metrics.
	DefaultEnable = false

	// DefaultEnable is the default setting for enablingPrometheus metric mapping.
	DefaultPrometheus = false

	// DefaultInterval holds the default metrics collection interval.
	DefaultInterval = 10 * time.Second
)

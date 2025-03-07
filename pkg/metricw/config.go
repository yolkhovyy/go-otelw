package metricw

import (
	"time"

	"github.com/yolkhovyy/go-otelw/pkg/collector"
)

type Config struct {
	Enable     bool             `yaml:"enable" mapstructure:"enable"`
	Prometheus bool             `yaml:"prometheus" mapstructure:"prometheus"`
	Interval   time.Duration    `yaml:"interval" mapstructure:"interval"`
	Collector  collector.Config `yaml:"collector" mapstructure:"collector"`
}

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
	DefaultEnable     = false
	DefaultPrometheus = false
	DefaultInterval   = 10 * time.Second
)

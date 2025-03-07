package slogw

import (
	"time"

	"github.com/yolkhovyy/go-otelw/pkg/collector"
)

type Config struct {
	Enable     bool             `yaml:"enable" mapstructure:"enable"`
	Caller     bool             `yaml:"caller" mapstructure:"Caller"`
	Format     Format           `yaml:"format" mapstructure:"Format"`
	Level      string           `yaml:"level" mapstructure:"Level"`
	TimeFormat string           `yaml:"timeFormat" mapstructure:"TimeFormat"`
	Collector  collector.Config `yaml:"collector" mapstructure:"Collector"`
}

func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable
	defaults["Caller"] = DefaultCaller
	defaults["Format"] = DefaultFormat
	defaults["Level"] = DefaultLevel
	defaults["TimeFormat"] = DefaultTimeFormat

	for k, v := range collector.Defaults() {
		defaults["Collector."+k] = v
	}

	return defaults
}

const (
	DefaultEnable     = false
	DefaultCaller     = false
	DefaultFormat     = FormatJSON
	DefaultLevel      = "info"
	DefaultTimeFormat = time.RFC3339
)

package otelw

import (
	"github.com/yolkhovyy/go-otelw/pkg/metricw"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
)

type Config struct {
	Logger slogw.Config   `yaml:"logger" mapstructure:"Logger"`
	Tracer tracew.Config  `yaml:"tracer" mapstructure:"Tracer"`
	Metric metricw.Config `yaml:"metric" mapstructure:"Metric"`
}

func Defaults() map[string]any {
	defaults := make(map[string]any)

	for k, v := range slogw.Defaults() {
		defaults["Logger."+k] = v
	}

	for k, v := range tracew.Defaults() {
		defaults["Tracer."+k] = v
	}

	for k, v := range metricw.Defaults() {
		defaults["Metric."+k] = v
	}

	return defaults
}

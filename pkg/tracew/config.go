package tracew

import "github.com/yolkhovyy/go-otelw/pkg/collector"

type Config struct {
	Enable    bool             `yaml:"enable" mapstructure:"enable"`
	Collector collector.Config `yaml:"collector" mapstructure:"collector"`
}

func Defaults() map[string]any {
	defaults := make(map[string]any)

	defaults["Enable"] = DefaultEnable

	for k, v := range collector.Defaults() {
		defaults["Collector."+k] = v
	}

	return defaults
}

const DefaultEnable = false

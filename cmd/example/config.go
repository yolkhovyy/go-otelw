package main

import (
	"fmt"

	httpserver "github.com/yolkhovyy/go-otelw/cmd/example/internal/server/http"
	"github.com/yolkhovyy/go-otelw/otelw"
	"github.com/yolkhovyy/go-utilities/viperx"
)

type Config struct {
	otelw.Config `yaml:",inline" mapstructure:",squash"`
	HTTP         httpserver.Config `yaml:"http" mapstructure:"HTTP"`
}

func (c *Config) Load(configFile string, prefix string) error {
	vprx := viperx.New(configFile, prefix, nil)

	vprx.SetDefaults(otelw.Defaults())
	vprx.SetDefaults(httpserver.Defaults())

	if err := vprx.Load(c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

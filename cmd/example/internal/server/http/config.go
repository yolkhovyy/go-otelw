package http

import "time"

type Config struct {
	Port              int           `yaml:"port" mapstructure:"Port"`
	ShutdownTimeout   time.Duration `yaml:"shutdownTimeout" mapstructure:"ShutdownTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" mapstructure:"ReadHeaderTimeout"`
}

package http

import "time"

type Config struct {
	Port              int           `yaml:"port" mapstructure:"Port"`
	ShutdownTimeout   time.Duration `yaml:"shutdownTimeout" mapstructure:"ShutdownTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" mapstructure:"ReadHeaderTimeout"`
}

func Defaults() map[string]any {
	return map[string]any{
		"HTTP.Port":              DefaultPort,
		"HTTP.ShutdownTimeout":   DefaultShutdownTimeout,
		"HTTP.ReadHeaderTimeout": DefaultReadHeaderTimeout,
	}
}

const (
	DefaultPort              = 8080
	DefaultShutdownTimeout   = 5 * time.Second
	DefaultReadHeaderTimeout = 1 * time.Second
)

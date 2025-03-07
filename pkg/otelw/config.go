package otelw

import (
	"time"
)

type Config struct {
	Logger Logger `yaml:"logger" mapstructure:"Logger"`
	Tracer Tracer `yaml:"tracer" mapstructure:"Tracer"`
	Metric Metric `yaml:"metric" mapstructure:"Metric"`
}

type Logger struct {
	Enable     bool      `yaml:"enable" mapstructure:"enable"`
	Caller     bool      `yaml:"caller" mapstructure:"Caller"`
	Format     Format    `yaml:"format" mapstructure:"Format"`
	Level      string    `yaml:"level" mapstructure:"Level"`
	Collector  Collector `yaml:"collector" mapstructure:"Collector"`
	TimeFormat string    `yaml:"timeFormat" mapstructure:"TimeFormat"`
}

type Tracer struct {
	Enable    bool      `yaml:"enable" mapstructure:"enable"`
	Collector Collector `yaml:"collector" mapstructure:"collector"`
}

type Metric struct {
	Enable     bool          `yaml:"enable" mapstructure:"enable"`
	Prometheus bool          `yaml:"prometheus" mapstructure:"prometheus"`
	Interval   time.Duration `yaml:"interval" mapstructure:"interval"`
	Collector  Collector     `yaml:"collector" mapstructure:"collector"`
}

func Defaults() map[string]any {
	return map[string]any{
		"Logger.Caller":               DefaultLoggerCaller,
		"Logger.Format":               DefaultLoggerFormat,
		"Logger.Level":                DefaultLoggerLevel,
		"Logger.TimeFormat":           DefaultLoggerTimeFormat,
		"Logger.Collector.Protocol":   DefaultLoggerCollectorProtocol,
		"Logger.Collector.Connection": DefaultLoggerCollectorConnection,

		"Tracer.Enable":               DefaultTracerEnable,
		"Tracer.Collector.Protocol":   DefaultTracerCollectorProtocol,
		"Tracer.Collector.Connection": DefaultTracerCollectorConnection,

		"Metric.Enable":               DefaultMetricEnable,
		"Metric.Prometheus":           DefaultMetricPrometheus,
		"Metric.Interval":             DefaultMetricInterval,
		"Metric.Collector.Protocol":   DefaultMetricCollectorProtocol,
		"Metric.Collector.Connection": DefaultMetricCollectorConnection,
	}
}

const (
	DefaultLoggerCaller              = false
	DefaultLoggerFormat              = FormatJSON
	DefaultLoggerLevel               = "info"
	DefaultLoggerTimeFormat          = time.RFC3339
	DefaultLoggerCollectorProtocol   = GRPC
	DefaultLoggerCollectorConnection = "localhost:4317"

	DefaultTracerEnable              = false
	DefaultTracerCollectorProtocol   = GRPC
	DefaultTracerCollectorConnection = "localhost:4317"

	DefaultMetricEnable              = false
	DefaultMetricPrometheus          = false
	DefaultMetricInterval            = 10 * time.Second
	DefaultMetricCollectorProtocol   = GRPC
	DefaultMetricCollectorConnection = "localhost:4317"
)

package otelw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yolkhovyy/go-utilities/viperx"
)

//nolint:funlen
func TestBaseLoad(t *testing.T) {
	t.Parallel()

	type args struct {
		configFile string
	}

	type want struct {
		err    bool
		config Config
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid",
			args: args{
				configFile: "test_data/valid_config.yml",
			},
			want: want{
				err: false,
				config: Config{
					Logger: Logger{
						Caller:     true,
						Format:     FormatJSON,
						Level:      "trace",
						TimeFormat: time.RFC3339Nano,
						Collector: Collector{
							Protocol:   GRPC,
							Connection: "foo:4242",
						},
					},
					Tracer: Tracer{
						Enable: true,
						Collector: Collector{
							Protocol:   GRPC,
							Connection: "foo:4242",
						},
					},
					Metric: Metric{
						Enable:     true,
						Prometheus: true,
						Interval:   42 * time.Second,
						Collector: Collector{
							Protocol:   GRPC,
							Connection: "foo:4242",
						},
					},
				},
			},
		},
		// {
		// 	name: "default",
		// 	args: args{
		// 		configFile: "test_data/default_config.yml",
		// 	},
		// 	want: want{
		// 		err: false,
		// 		config: Config{
		// 			Logger: Logger{
		// 				Caller:     DefaultLoggerCaller,
		// 				Format:     DefaultLoggerFormat,
		// 				Level:      DefaultLoggerLevel,
		// 				TimeFormat: DefaultLoggerTimeFormat,
		// 				Collector: Collector{
		// 					Protocol:   DefaultLoggerCollectorProtocol,
		// 					Connection: DefaultLoggerCollectorConnection,
		// 				},
		// 			},
		// 			Tracer: Tracer{
		// 				Enable: DefaultTracerEnable,
		// 				Collector: Collector{
		// 					Protocol:   DefaultTracerCollectorProtocol,
		// 					Connection: DefaultTracerCollectorConnection,
		// 				},
		// 			},
		// 			Metric: Metric{
		// 				Enable:     DefaultMetricEnable,
		// 				Prometheus: DefaultMetricPrometheus,
		// 				Interval:   DefaultMetricInterval,
		// 				Collector: Collector{
		// 					Protocol:   DefaultMetricCollectorProtocol,
		// 					Connection: DefaultMetricCollectorConnection,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "invalid",
		// 	args: args{
		// 		configFile: "test_data/invalid_config.yml",
		// 	},
		// 	want: want{
		// 		err: true,
		// 	},
		// },
		// {
		// 	name: "non-existing",
		// 	args: args{
		// 		configFile: "test_data/non-existing_config.yml",
		// 	},
		// 	want: want{
		// 		err: true,
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vprx := viperx.New(test.args.configFile, "FOO", nil)
			vprx.SetDefaults(Defaults())

			config := Config{}
			err := vprx.Load(&config)

			if test.want.err {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.want.config, config)
			}
		})
	}
}

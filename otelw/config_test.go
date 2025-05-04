package otelw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yolkhovyy/go-otelw/otelw/metricw"
	"github.com/yolkhovyy/go-otelw/otelw/otlp"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"github.com/yolkhovyy/go-otelw/otelw/tracew"
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
					Logger: slogw.Config{
						Caller:     true,
						Format:     slogw.JSON,
						Level:      "trace",
						TimeFormat: time.RFC3339Nano,
						OTLP: otlp.Config{
							Protocol: otlp.GRPC,
							Endpoint: "foo:4242",
						},
					},
					Tracer: tracew.Config{
						Enable: true,
						OTLP: otlp.Config{
							Protocol: otlp.GRPC,
							Endpoint: "foo:4242",
						},
					},
					Metric: metricw.Config{
						Enable:     true,
						Prometheus: true,
						Interval:   42 * time.Second,
						OTLP: otlp.Config{
							Protocol: otlp.GRPC,
							Endpoint: "foo:4242",
						},
					},
				},
			},
		},
		{
			name: "default",
			args: args{
				configFile: "test_data/default_config.yml",
			},
			want: want{
				err: false,
				config: Config{
					Logger: slogw.Config{
						Caller:     slogw.DefaultCaller,
						Format:     slogw.DefaultFormat,
						Level:      slogw.DefaultLevel,
						TimeFormat: slogw.DefaultTimeFormat,
						OTLP: otlp.Config{
							Protocol: otlp.DefaultProtocol,
							Endpoint: otlp.DefaultEndpoint,
						},
					},
					Tracer: tracew.Config{
						Enable: tracew.DefaultEnable,
						OTLP: otlp.Config{
							Protocol: otlp.DefaultProtocol,
							Endpoint: otlp.DefaultEndpoint,
						},
					},
					Metric: metricw.Config{
						Enable:     metricw.DefaultEnable,
						Prometheus: metricw.DefaultPrometheus,
						Interval:   metricw.DefaultInterval,
						OTLP: otlp.Config{
							Protocol: otlp.DefaultProtocol,
							Endpoint: otlp.DefaultEndpoint,
						},
					},
				},
			},
		},
		{
			name: "invalid",
			args: args{
				configFile: "test_data/invalid_config.yml",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "non-existing",
			args: args{
				configFile: "test_data/non-existing_config.yml",
			},
			want: want{
				err: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vprx := viperx.New(test.args.configFile, "FOO", nil)
			vprx.SetDefaults(Defaults())

			config := Config{}
			err := vprx.Load(&config)

			if test.want.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.want.config, config)
			}
		})
	}
}

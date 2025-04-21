package slogw

import (
	"context"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yolkhovyy/go-otelw/test"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

//nolint:funlen
func TestSlog(t *testing.T) {
	t.Parallel()

	type args struct {
		config     Config
		attributes []attribute.KeyValue
		writers    []io.Writer
		message    string
		traceID    trace.TraceID
		spanID     trace.SpanID
	}

	type want struct {
		Type  any
		RxLog *regexp.Regexp
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "slogw",
			args: args{
				config: Config{
					Enable: true,
					Level:  "trace",
				},
				attributes: []attribute.KeyValue{
					semconv.ServiceNameKey.String("slogw"),
					semconv.ServiceVersionKey.String("v0.1.0"),
				},
				writers: []io.Writer{},
				message: "test log",
				traceID: trace.TraceID{
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
					0x18, 0x19, 0x1a, 0x1b,
					0x1c, 0x1d, 0x1e, 0x1f,
				},
				spanID: trace.SpanID{
					0x20, 0x21, 0x22, 0x23,
					0x24, 0x25, 0x26, 0x27,
				},
			},
			want: want{
				Type: &Logger{},
				//nolint:lll
				RxLog: regexp.MustCompile(
					`{"Timestamp":"` + test.RxTime + `","ObservedTimestamp":"` + test.RxTime + `","Severity":\d{1,2},"SeverityText":"(?i)\b(trace|debug|info|warn|error|fatal|panic)\d{0,2}\b","Body":{"Type":"String","Value":"test log"},"Attributes":\[\],"TraceID":"101112131415161718191a1b1c1d1e1f","SpanID":"2021222324252627","TraceFlags":"00","Resource":\[{"Key":"service\.name","Value":{"Type":"STRING","Value":"slogw"}},{"Key":"service\.version","Value":{"Type":"STRING","Value":"v\d+\.\d+\.\d+"}},{"Key":"telemetry\.sdk\.language","Value":{"Type":"STRING","Value":"go"}},{"Key":"telemetry\.sdk\.name","Value":{"Type":"STRING","Value":"opentelemetry"}},{"Key":"telemetry\.sdk\.version","Value":{"Type":"STRING","Value":"` + test.RxTelemetrySDKVersion + `"}}\],"Scope":{"Name":"slogw","Version":"","SchemaURL":"","Attributes":{}},"DroppedAttributes":0}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var builder strings.Builder
			test.args.writers = append(test.args.writers, &builder)

			spanContext := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: test.args.traceID,
				SpanID:  test.args.spanID,
			})
			ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

			logger, err := Configure(ctx, test.args.config, test.args.attributes, test.args.writers...)
			require.NoError(t, err)
			assert.IsType(t, test.want.Type, logger)

			logger.InfoContext(ctx, test.args.message)

			err = logger.ForceFlush(ctx)
			require.NoError(t, err)

			output := builder.String()
			if test.args.config.Format == Console {
				output = removeColorFormatting(output)
			}

			assert.Regexp(t, test.want.RxLog, output)

			builder.Reset()
			logger.DebugContext(ctx, test.args.message)

			err = logger.ForceFlush(ctx)
			require.NoError(t, err)

			output = builder.String()
			if test.args.config.Format == Console {
				output = removeColorFormatting(output)
			}

			assert.Regexp(t, test.want.RxLog, output)
		})
	}
}

func removeColorFormatting(input string) string {
	re := regexp.MustCompile(`\033\[[0-9;]*[mK]`)

	return re.ReplaceAllString(input, "")
}

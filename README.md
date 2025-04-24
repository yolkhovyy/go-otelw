# OpenTelemetry Toolkit for Golang  

![License](https://img.shields.io/github/license/yolkhovyy/go-otelw)
![GitHub Tag](https://img.shields.io/github/v/tag/yolkhovyy/go-otelw)
[![Go Reference](https://pkg.go.dev/badge/github.com/yolkhovyy/go-otelw.svg)](https://pkg.go.dev/github.com/yolkhovyy/go-otelw)
[![Go Report Card](https://goreportcard.com/badge/github.com/yolkhovyy/go-otelw)](https://goreportcard.com/report/github.com/yolkhovyy/go-otelw)

🚀 **OpenTelemetry made easy for Golang**  
✨ **The one-stop place for Golang & OpenTelemetry**  
🌀 **Helps to wrap your head around OpenTelemetry**  

Lightweight OpenTelemetry toolkit for Go, with plug-and-play examples for many observability platforms.
Simplifies setup by wrapping OpenTelemetry with `Configure` and `Shutdown` utility functions — hence the name `go-otelw` which is pronounced /ˈɡuːtldʌb/, short for Go OpenTelemetry Wrapper.

Observability backend examples included:
  * [Datadog](docs/datadog.md)
  * [Dynatrace](docs/dynatrace.md)
  * [Elasticsearch Kibana](docs/elasticsearch-kibana.md)
  * [Grafana Cloud](docs/grafana-cloud.md)
  * [Grafana Cloud Alloy](./docs/grafana-cloud-alloy.md)
  * [Grafana Loki, Jaeger, Prometheus](./docs/grafana-loki-jaeger-prometheus.md)
  * [Grafana Loki, Tempo, Prometheus](./docs/grafana-loki-tempo-prometheus.md)
  * [Honeycomb](./docs/honeycomb.md)
  * [New Relic](./docs/new-relic.md)
  * [OpenObserve](./docs/openobserve.md)
  * [Uptrace](./docs/uptrace.md)

⚠️ This project is pre-v1.0.0 and may change.

## Overview
The diagram below illustrates how telemetry from the Example service, instrumented with `go-otelw`, can be routed to any OTEL-compatible backend—or to multiple backends at once. Simply instrument your service with `go-otelw` and use the provided OTEL Collector configuration examples to route telemetry to the observability backend of your choice.

![Overview](./docs/diagrams/overview.png)

## Content
* [How to Integrate OpenTelemetry with go-otelw](#how-to-integrate-opentelemetry-with-go-otelw)
  * [Config Types](#config-types)
  * [Configure and Shutdown Utility Functions](#configure-and-shutdown-utility-functions)
  * [Tracing and Logging Examples](#tracing-and-logging-examples)
* Build and Run the Examples
  * [Datadog](docs/datadog.md)
  * [Dynatrace](docs/dynatrace.md)
  * [Elasticsearch Kibana](docs/elasticsearch-kibana.md)
  * [Grafana Cloud](docs/grafana-cloud.md)
  * [Grafana Cloud Alloy](./docs/grafana-cloud-alloy.md)
  * [Grafana Loki, Jaeger, Prometheus](./docs/grafana-loki-jaeger-prometheus.md)
  * [Grafana Loki, Tempo, Prometheus](./docs/grafana-loki-tempo-prometheus.md)
  * [Honeycomb](./docs/honeycomb.md)
  * [New Relic](./docs/new-relic.md)
  * [OpenObserve](./docs/openobserve.md)
  * [Uptrace](./docs/uptrace.md)

## Package Content
* The OpenTelemetry [Wrapper](./otelw/) itself
* Usage [Example](./cmd/example/) - HTTP Echo Service
* Docker [Compose](./docker-compose.yml) to run the Echo Service and its dependencies
* [Config](./config/) files for 3rd-party dependencies
	* [Grafana](./config/grafana/)
	* [Jaeger](./config/jaeger/)
	* [Loki](./config/loki/)
	* [Prometheus](./config/prometheus/)
	* [Promtail](./config/promtail/)
	* [Tempo](./config/tempo/)
	* [Uptrace](./config/uptrace/)

## How to Integrate OpenTelemetry with go-otelw
### Install
```bash
go get github.com/yolkhovyy/go-otelw@latest
```

### Config Types
`go-otelw` provides convenience config types for logger, tracer, metric in [otelw/config.go](https://github.com/yolkhovyy/go-otelw/blob/main/otelw/config.go#L11-L15)
```go
type Config struct {
	// Logging configuration
	Logger slogw.Config   `yaml:"logger" mapstructure:"Logger"`
	
	// Tracing configuration
	Tracer tracew.Config  `yaml:"tracer" mapstructure:"Tracer"`
	
	// Metrics configuration
	Metric metricw.Config `yaml:"metric" mapstructure:"Metric"`
}
```
and OTEL Collector in [otelw/collector/config.go](https://github.com/yolkhovyy/go-otelw/blob/main/otelw/collector/config.go#L6-L18)
```go
type Config struct {
	// Protocol to use for telemetry collection - GRPC (default), HTTP.
	Protocol Protocol `yaml:"protocol" mapstructure:"protocol"`

	// Address of the telemetry collector service
	Connection string `yaml:"connection" mapstructure:"connection"`

	// Whether to use an insecure connection (without TLS)
	Insecure bool `yaml:"insecure" mapstructure:"insecure"`

	// TLS configuration settings
	TLS TLS `yaml:"tls" mapstructure:"tls"`
}
```
`go-otelw` configuration can be loaded from yaml or json files on application startup. An example of yaml configuration is in [cmd/example/config.yml](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/config.yml)

### Configure and Shutdown Utility Functions
`go-otelw` simplifies the use of OpenTelemetry by providing `Configure` and `Shutdown` utility functions for logger, tracer and metric.

#### All-in Configuration and Shutdown
See [cmd/example/main.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/main.go#L60-L75)

```golang
	serviceAttributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version.Tag),
	}
	
	logger, tracer, metric, err := otelw.Configure(ctx, config.Config, serviceAttributes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otelw configure: %v", err)

		return osx.ExitFailure
	}

	defer func() {
		err := errors.Join(err,
			metric.Shutdown(ctx),
			tracer.Shutdown(ctx),
			logger.Shutdown(ctx))
		if err != nil {
			fmt.Fprintf(os.Stderr, "otelw shutdown: %v", err)
		}
	}()
```

#### Individual Logger, Tracer and Metric Configuration
See [cmd/example/internal/otelw/otelw.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/internal/otelw/otelw.go#L21-L34)

```golang
	logger, err := slogw.Configure(ctx, config.Logger, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("slogw configure: %w", err)
	}

	tracer, err := tracew.Configure(ctx, config.Tracer, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("tracew configure: %w", err)
	}

	metric, err := metricw.Configure(ctx, config.Metric, attrs, writers...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("metricw configure: %w", err)
	}

```

### Tracing and Logging Examples
See [cmd/example/internal/domain/domain.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/internal/domain/domain.go#L75-L110)

```golang
	ctx, span := tracew.Start(ctx, "echo", "worker"+strconv.Itoa(sequence))
	defer func() { span.End(err) }()

	logger := slogw.DefaultLogger()
	logger.InfoContext(ctx, msg, 
		slog.Int("sequence", sequence),
		slog.String("input", input),
	)
```

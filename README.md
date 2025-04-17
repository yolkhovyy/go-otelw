# OpenTelemetry Wrapper and Examples for Golang  

🚀 **OpenTelemetry made easy for Golang**  
✨ **The one-stop place for Golang & OpenTelemetry**  

This is an OpenTelemetry Wrapper and Examples for Golang. Its goal is to simplify OpenTelemetry integration and usage in Golang.

Pronounced /ˈɡuːtldʌb/

Examples included:
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

## Overview
By modifying the OTEL Collector's config.yml, the telemetry flow can be directed to any OTEL-compatible backend—or even to multiple backends simultaneously. The diagram below illustrates this process.

![Overview](./docs/diagrams/overview.png)

## Content
* How to Integrate OpenTelemetry with go-otelw
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
* The [wrapper](./pkg/) itself
* Usage [example](./cmd/example/) - HTTP Echo Service
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
### Config Types
`go-otelw` provides convenience config types for logger, tracer, metric in [pkg/otelw/config.go](https://github.com/yolkhovyy/go-otelw/blob/main/pkg/otelw/config.go#L11-L15)
```go
type Config struct {
	Logger slogw.Config   `yaml:"logger" mapstructure:"Logger"` // Logging configuration
	Tracer tracew.Config  `yaml:"tracer" mapstructure:"Tracer"` // Tracing configuration
	Metric metricw.Config `yaml:"metric" mapstructure:"Metric"` // Metrics configuration
}
```
and OTEL Collector in [pkg/collector/config.go](https://github.com/yolkhovyy/go-otelw/blob/main/pkg/collector/config.go#L6-L18)
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
`go-otelw` simplifies using of OpenTelemetry by providing `Configure` and `Shutdown` utility functions for logger, tracer and metric.

See usage example in [cmd/example/main.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/main.go#L60-L75)

```golang
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

# OpenTelemetry Wrapper for Golang  

🚀 **OpenTelemetry made easy for Golang**  
✨ **The one-stop place for Golang & OpenTelemetry**  

This is an OpenTelemetry Wrapper for Golang. Its goal is to simplify OpenTelemetry integration and usage in Golang.

Pronounced /ˈɡuːtldʌb/

Examples included:
  * Datadog
  * Dynatrace
  * Elasticsearch, Kibana
  * Grafana Loki, Jaeger/Tempo, Prometheus
  * Honeycomb
  * New Relic
  * OpenObserve
  * Uptrace

## Overview
![Overview](./docs/diagrams/overview.png)

## Content
* How to Integrate OpenTelemetry with go-otelw
  * [Configuration and Shutdown](#configuration-and-shutdown)
  * [Logger and Tracer Example](#logger-and-tracer-example)
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

### Configuration and Shutdown

See [cmd/example/main.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/main.go#L60-L75)

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

### Logger and Tracer Example

See [cmd/example/internal/domain/domain.go](https://github.com/yolkhovyy/go-otelw/blob/main/cmd/example/internal/domain/domain.go#L75-L110)

```golang
	ctx, span := tracew.Start(ctx, "echo", "worker"+strconv.Itoa(sequence))
	defer func() { span.End(err) }()

	logger := slogw.DefaultLogger()
```

## Miscellaneous

**OpenTelemetry SDK Golang**
Examples:
* [Product Catalog](https://opentelemetry.io/docs/demo/services/product-catalog/)
* [Checkout](https://opentelemetry.io/docs/demo/services/checkout/)

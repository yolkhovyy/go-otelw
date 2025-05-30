# Grafana Loki, Jaeger, Prometheus Integration
![Grafana Loki Jaeger](./diagrams/grafana-loki-jaeger.png)

**Build and run the Example:**
```bash
make doco-build-up
```

This will start the `Example` Echo Service, and the telemetry services - `OTEL collector`, `Grafana/Loki`, `Jaeger`, and `Prometheus`.

**Make a few HTTP requests to the Example Echo Service:**
```bash
./test/scripts/echo.sh
./test/scripts/echo.sh hey 10
```

**Observe logs, traces and metrics in OTEL Collector's logs:**
```bash
docker compose logs -f otel-collector
```

**Explore logs, traces and metrics:**
* Logs
  * Garfana: http://localhost:3000
    * Explore - Select data source `Loki` - Select label `service_name` - Select value `example` - Click `Run query`
* Traces
  * Jaeger: http://localhost:16686
  * Garfana: http://localhost:3000
    * Explore - Select data source `Jaeger` - Click Query type `Search` - Select service name `example` - Click `Run query` 
* Metrics
  * Prometheus: http://localhost:9090
  * Garfana: http://localhost:3000
    * Explore - Select Prometheus data source - Select a metric - Click `Run query`

**Stop the services:**
```bash
make doco-down
```

## Miscellaneous

**Grafana**
* [Jaeger provisioning](https://grafana.com/docs/grafana/latest/datasources/jaeger/#provision-the-data-source)
* [Tempo provisioning](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/#example-file)
* Useful Grafana dashboards:
  * [Logs/App, ID: 13639](https://grafana.com/grafana/dashboards/13639-logs-app/)

# TODO

## Tracer Provider

Consider for production:
* [OTEL example](https://opentelemetry.io/docs/demo/services/checkout/)

```golang
    provider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
    )
```

# TODO

## Tracer Provider

Consider for production:
* [OTEL example](https://opentelemetry.io/docs/demo/services/checkout/)

```golang
    provider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
    )
```

## Backends

* Finalize splunk
* [dash0](https://www.dash0.com/)

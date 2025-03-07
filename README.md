# Go OpenTelemetry Wrapper

This is a Go OpenTelemetry playground project. It provides a wrapper for OpenTelemetry with the goal of simplifying its integration and usage.

Pronounced as /ˈɡuːtldʌb/


## Content
* The [wraper](./pkg/) itself
* Usage [example](./cmd/example/) - HTTP Echo Service
* Docker [Compose](./docker-compose.yml) to run the Echo Service and its dependencies
* [Configuration](./config/) files for 3rd-party dependencies

## TODO
* Rethink otelw.Configure() return - is it possible to do Shutdown() using only globals?

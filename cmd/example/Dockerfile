FROM golang:1.23.6-alpine3.20 AS builder
ARG APP_VERSION
ARG LDFLAGS="-X github.com/yolkhovyy/go-otelw/cmd/example/version.Tag=${APP_VERSION}"

WORKDIR /app
COPY ./cmd/example ./cmd/example
COPY ./pkg ./pkg
COPY go.mod go.sum ./
RUN go mod download && \
    go build -ldflags="${LDFLAGS}" -o main ./cmd/example

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY ./cmd/example/config.yml .
EXPOSE 8080
CMD ["./main", "--config", "config.yml"]


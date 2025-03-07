package otelw

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/credentials"
)

func FindAttribute(attrs []attribute.KeyValue, key attribute.Key) string {
	for _, attr := range attrs {
		if attr.Key == key && attr.Valid() {
			return attr.Value.AsString()
		}
	}

	return ""
}

//nolint:ireturn
func TLSCredentials(config Collector) (credentials.TransportCredentials, error) {
	tlsConfig, err := TLSConfig(config)
	if err != nil {
		return nil, fmt.Errorf("tls config: %w", err)
	}

	tlsCreds := credentials.NewTLS(tlsConfig)

	return tlsCreds, nil
}

func TLSConfig(config Collector) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.TLS.ClientCrt, config.TLS.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	caCert, err := os.ReadFile(config.TLS.CAFile)
	if err != nil {
		return nil, fmt.Errorf("load ca cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("append ca cert: %w", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		ServerName:         config.Connection,
	}

	return tlsConfig, nil
}

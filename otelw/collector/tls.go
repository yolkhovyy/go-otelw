package collector

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// TLSCredentials creates gRPC transport credentials using TLS settings from the provided Config.
// It loads the TLS configuration and returns transport credentials or an error if loading fails.
func TLSCredentials( //nolint:ireturn
	config Config,
) (credentials.TransportCredentials, error) {
	tlsConfig, err := TLSConfig(config)
	if err != nil {
		return nil, fmt.Errorf("tls config: %w", err)
	}

	tlsCreds := credentials.NewTLS(tlsConfig)

	return tlsCreds, nil
}

// TLSConfig generates a tls.Config from the provided Config structure.
// It loads the client certificate, key, and CA certificate to configure mutual TLS authentication.
func TLSConfig(config Config) (*tls.Config, error) {
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

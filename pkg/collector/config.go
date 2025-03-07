package collector

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

type Config struct {
	Protocol   Protocol `yaml:"protocol" mapstructure:"protocol"`
	Connection string   `yaml:"connection" mapstructure:"connection"`
	Insecure   bool     `yaml:"insecure" mapstructure:"insecure"`
	TLS        TLS      `yaml:"tls" mapstructure:"tls"`
}

type TLS struct {
	ClientCrt string `yaml:"clientCrt" mapstructure:"clientCrt"`
	ClientKey string `yaml:"clientKey" mapstructure:"clientKey"`
	CAFile    string `yaml:"caFile" mapstructure:"caFile"`
}

type Protocol string

const (
	GRPC Protocol = "grpc"
	HTTP Protocol = "http"
)

func (p Protocol) String() string {
	return string(p)
}

func Defaults() map[string]any {
	return map[string]any{
		"Protocol":   DefaultProtocol,
		"Connection": DefaultConnection,
	}
}

const (
	DefaultProtocol   = GRPC
	DefaultConnection = "localhost:4317"
)

//nolint:ireturn
func TLSCredentials(config Config) (credentials.TransportCredentials, error) {
	tlsConfig, err := TLSConfig(config)
	if err != nil {
		return nil, fmt.Errorf("tls config: %w", err)
	}

	tlsCreds := credentials.NewTLS(tlsConfig)

	return tlsCreds, nil
}

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

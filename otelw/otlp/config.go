package otlp

// Config represents the configuration settings for a OTEL protocol.
// It includes details about the protocol, endpoint settings, security options,
// and TLS configuration.
type Config struct {
	// Protocol to use for telemetry collection - GRPC (default), HTTP.
	Protocol Protocol `json:"protocol" yaml:"protocol" mapstructure:"protocol"`

	// Endpoint for OTLP protocol.
	Endpoint string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`

	// Whether to use an insecure endpoint (without TLS).
	Insecure bool `json:"insecure" yaml:"insecure" mapstructure:"insecure"`

	// Path to the client certificate file.
	ClientCertificate string `json:"clientCertificate" yaml:"clientCertificate" mapstructure:"clientCertificate"`

	// Path to the client private key file.
	ClientKey string `json:"clientKey" yaml:"clientKey" mapstructure:"clientKey"`

	// Path to the Certificate Authority (CA) file.
	Certificate string `json:"certificate" yaml:"certificate" mapstructure:"certificate"`
}

// Defaults returns a map of default configuration values for the otlp package.
// It sets a default protocol and endpoint address.
func Defaults() map[string]any {
	return map[string]any{
		"Protocol": DefaultProtocol,
		"Endpoint": DefaultEndpoint,
	}
}

const (
	// Default protocol for telemetry collection (gRPC).
	DefaultProtocol = GRPC

	// Default otlp endpoint.
	DefaultEndpoint = "localhost:4317"
)

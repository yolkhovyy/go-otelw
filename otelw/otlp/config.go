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
	ClientCertificate string `json:"client_certificate" yaml:"clientCertificate" mapstructure:"clientCertificate"`

	// Path to the client private key file.
	ClientKey string `json:"client_key" yaml:"clientKey" mapstructure:"clientKey"`

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
	// DefaultProtocol for telemetry collection (gRPC).
	DefaultProtocol = GRPC

	// DefaultEndpoint for the OTLP endpoint.
	DefaultEndpoint = "localhost:4317"
)

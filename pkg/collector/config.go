package collector

// Config represents the configuration settings for a telemetry collector.
// It includes details about the protocol, connection settings, security options,
// and TLS configuration.
type Config struct {
	// Protocol to use for telemetry collection - GRPC (default), HTTP.
	Protocol Protocol `yaml:"protocol" mapstructure:"protocol"`

	// Address of the telemetry collector service
	Connection string `yaml:"connection" mapstructure:"connection"`

	// Whether to use an insecure connection (without TLS)
	Insecure bool `yaml:"insecure" mapstructure:"insecure"`

	// TLS configuration settings
	TLS TLS `yaml:"tls" mapstructure:"tls"`
}

// TLS holds the TLS security configuration options for secure telemetry connections.
type TLS struct {
	// Path to the client certificate file
	ClientCrt string `yaml:"clientCrt" mapstructure:"clientCrt"`

	// Path to the client private key file
	ClientKey string `yaml:"clientKey" mapstructure:"clientKey"`

	// Path to the Certificate Authority (CA) file
	CAFile string `yaml:"caFile" mapstructure:"caFile"`
}

// Defaults returns a map of default configuration values for the collector package.
// It sets a default protocol and connection address.
func Defaults() map[string]any {
	return map[string]any{
		"Protocol":   DefaultProtocol,
		"Connection": DefaultConnection,
	}
}

const (
	// Default protocol for telemetry collection (gRPC).
	DefaultProtocol = GRPC

	// Default collector service address.
	DefaultConnection = "localhost:4317"
)

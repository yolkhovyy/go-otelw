package otelw

type Collector struct {
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

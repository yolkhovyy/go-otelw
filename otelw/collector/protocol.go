package collector

// Protocol defines a type for supported communication protocols.
type Protocol string

const (
	// GRPC represents the gRPC protocol.
	GRPC Protocol = "grpc"
	// HTTP represents the HTTP protocol.
	HTTP Protocol = "http"
)

// String returns the string representation of the Protocol.
func (p Protocol) String() string {
	return string(p)
}

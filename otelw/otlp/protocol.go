package otlp

// Protocol defines a type for supported communication protocols.
type Protocol string

const (
	// GRPC represents the gRPC protocol.
	GRPC Protocol = "grpc"
	// HTTP represents the HTTP Protobuf protocol.
	HTTP Protocol = "http/protobuf"
)

// String returns the string representation of the Protocol.
func (p Protocol) String() string {
	return string(p)
}

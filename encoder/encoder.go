package encoder

// Encoder defines the interface for types that handle encoding and decoding of data.
type Encoder interface {
	// Marshal encodes a Go value into a byte slice (e.g., JSON, CBOR, etc.).
	Marshal(value any) ([]byte, error)
	// Unmarshal decodes a byte slice into a Go value (e.g., JSON, CBOR, etc.).
	Unmarshal(data []byte, value any) error
}

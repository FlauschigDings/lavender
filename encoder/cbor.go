package encoder

import (
	"github.com/fxamacker/cbor"
)

// CBorEncoder is a concrete implementation of the Encoder interface that uses
// the CBOR (Concise Binary Object Representation) encoding format for serializing
// and deserializing data.
type CBorEncoder struct {
	EncOpts cbor.EncOptions
}

// CBorEncoder implements the Encoder interface by providing methods
// for marshaling and unmarshaling data in the CBOR format.
var _ Encoder = new(CBorEncoder)

// NewCBorEncoder creates and returns a new instance of CBorEncoder with default
// encoding options (unsorted keys).
func NewCBorEncoder() *CBorEncoder {
	return NewCBorCustomEncoder(
		cbor.PreferredUnsortedEncOptions(),
	)
}

// NewCBorCustomEncoder creates and returns a new instance of CBorEncoder with custom
// encoding options passed as an argument.
func NewCBorCustomEncoder(encOpts cbor.EncOptions) *CBorEncoder {
	return &CBorEncoder{
		EncOpts: encOpts,
	}
}

// Marshal implements Encoder. It converts a Go value into a byte slice using the CBOR
// encoding format with the specified encoding options.
func (g *CBorEncoder) Marshal(value any) ([]byte, error) {
	return cbor.Marshal(value, g.EncOpts)
}

// Unmarshal implements Encoder. It decodes a byte slice into a Go value using the CBOR
// encoding format.
func (g *CBorEncoder) Unmarshal(data []byte, value any) error {
	return cbor.Unmarshal(data, value)
}

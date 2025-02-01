package encoder

import "encoding/json"

// JsonEncoder implements the Encoder interface for encoding and decoding data in JSON format.
type JsonEncoder struct{}

// Ensure that JsonEncoder implements the Encoder interface.
var _ Encoder = new(JsonEncoder)

// NewJsonEncoder creates and returns a new instance of JsonEncoder.
func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

// Marshal implements the Encoder interface method for marshaling (encoding) a Go value into JSON bytes.
func (j *JsonEncoder) Marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

// Unmarshal implements the Encoder interface method for unmarshaling (decoding) JSON bytes into a Go value.
func (j *JsonEncoder) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

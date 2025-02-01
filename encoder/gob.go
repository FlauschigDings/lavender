package encoder

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// GobEncoder is a concrete implementation of the Encoder interface that uses
// the Gob encoding format for serializing and deserializing data.
type GobEncoder struct {
}

// GobEncoder implements the Encoder interface by providing methods
// for marshaling and unmarshaling data in the Gob format.
var _ Encoder = new(GobEncoder)

// NewGobEncoder creates and returns a new instance of GobEncoder.
func NewGobEncoder() *GobEncoder {
	return &GobEncoder{}
}

// Marshal implements Encoder. It converts a Go value into a byte slice using the Gob encoding format.
func (g *GobEncoder) Marshal(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	data, err := structToMap(value)
	if err != nil {
		return nil, err
	}
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal implements Encoder. It decodes a byte slice into a Go value using the Gob encoding format.
func (g *GobEncoder) Unmarshal(data []byte, value any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	ret := make(map[string]interface{})
	if err := dec.Decode(&ret); err != nil {
		return err
	}
	return mapToStruct(ret, value)
}

// structToMap converts a Go struct to a map[string]interface{} using JSON encoding/decoding.
func structToMap(s interface{}) (map[string]interface{}, error) {
	// Marshal struct to JSON
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to map
	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// mapToStruct converts a map[string]interface{} back into a Go struct using JSON encoding/decoding.
func mapToStruct(m map[string]interface{}, s interface{}) error {
	// Marshal map to JSON
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Unmarshal JSON to struct
	err = json.Unmarshal(bytes, s)
	if err != nil {
		return err
	}

	return nil
}

package json

import "encoding/json"

type JSONEncoder struct {
	// Empty
}

func New() *JSONEncoder {
	return &JSONEncoder{}
}

// Encode
func (js *JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode
func (js *JSONEncoder) Decode(data []byte, vPtr interface{}) error {
	return json.Unmarshal(data, vPtr)
}

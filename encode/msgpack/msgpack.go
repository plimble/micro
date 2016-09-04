package msgpack

import "github.com/tinylib/msgp/msgp"

type MsgPackEncoder struct {
	// Empty
}

func New() *MsgPackEncoder {
	return &MsgPackEncoder{}
}

// Encode
func (js *MsgPackEncoder) Encode(v msgp.Marshaler) ([]byte, error) {
	return v.MarshalMsg(nil)
}

// Decode
func (js *MsgPackEncoder) Decode(data []byte, vPtr msgp.Unmarshaler) error {
	_, err := vPtr.UnmarshalMsg(data)
	return err
}

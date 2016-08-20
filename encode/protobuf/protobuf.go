package protobuf

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// ProtobufEncoder is a protobuf implementation for EncodedConn
// This encoder will use the builtin protobuf lib to Marshal
// and Unmarshal structs.
type ProtobufEncoder struct {
	// Empty
}

func New() *ProtobufEncoder {
	return &ProtobufEncoder{}
}

var (
	ErrInvalidProtoMsgEncode = errors.New("nats: Invalid protobuf proto.Message object passed to encode")
	ErrInvalidProtoMsgDecode = errors.New("nats: Invalid protobuf proto.Message object passed to decode")
)

// Encode
func (pb *ProtobufEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	i, found := v.(proto.Message)
	if !found {
		return nil, ErrInvalidProtoMsgEncode
	}

	b, err := proto.Marshal(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Decode
func (pb *ProtobufEncoder) Decode(data []byte, vPtr interface{}) error {
	if _, ok := vPtr.(*interface{}); ok {
		return nil
	}
	i, found := vPtr.(proto.Message)
	if !found {
		return ErrInvalidProtoMsgDecode
	}

	err := proto.Unmarshal(data, i)
	if err != nil {
		return err
	}
	return nil
}

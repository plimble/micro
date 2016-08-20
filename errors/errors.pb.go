// Code generated by protoc-gen-go.
// source: github.com/plimble/micro/errors/errors.proto
// DO NOT EDIT!

/*
Package errors is a generated protocol buffer package.

It is generated from these files:
	github.com/plimble/micro/errors/errors.proto

It has these top-level messages:
	Errors
*/
package errors

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Errors struct {
	Code    int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty" gorethink:"code"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty" gorethink:"message"`
}

func (m *Errors) Reset()                    { *m = Errors{} }
func (m *Errors) String() string            { return proto.CompactTextString(m) }
func (*Errors) ProtoMessage()               {}
func (*Errors) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Errors)(nil), "errors.Errors")
}

func init() { proto.RegisterFile("github.com/plimble/micro/errors/errors.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 116 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xd2, 0x49, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0xc8, 0xc9, 0xcc, 0x4d, 0xca, 0x49, 0xd5, 0xcf,
	0xcd, 0x4c, 0x2e, 0xca, 0xd7, 0x4f, 0x2d, 0x2a, 0xca, 0x2f, 0x2a, 0x86, 0x52, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0x92, 0x19, 0x17, 0x9b, 0x2b, 0x98, 0x25, 0x24, 0xc4,
	0xc5, 0x92, 0x9c, 0x9f, 0x92, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1a, 0x04, 0x66, 0x0b, 0x49,
	0x70, 0xb1, 0xe7, 0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7, 0x4a, 0x30, 0x01, 0x85, 0x39, 0x83, 0x60,
	0xdc, 0x24, 0x36, 0xb0, 0x31, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x09, 0x0f, 0xe1, 0xd0,
	0x76, 0x00, 0x00, 0x00,
}

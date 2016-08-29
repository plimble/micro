package micro

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Header) DecodeMsg(dc *msgp.Reader) (err error) {
	var zajw uint32
	zajw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	if (*z) == nil && zajw > 0 {
		(*z) = make(Header, zajw)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zajw > 0 {
		zajw--
		var zbai string
		var zcmr interface{}
		zbai, err = dc.ReadString()
		if err != nil {
			return
		}
		zcmr, err = dc.ReadIntf()
		if err != nil {
			return
		}
		(*z)[zbai] = zcmr
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Header) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zwht, zhct := range z {
		err = en.WriteString(zwht)
		if err != nil {
			return
		}
		err = en.WriteIntf(zhct)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Header) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, uint32(len(z)))
	for zwht, zhct := range z {
		o = msgp.AppendString(o, zwht)
		o, err = msgp.AppendIntf(o, zhct)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Header) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zlqf uint32
	zlqf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	if (*z) == nil && zlqf > 0 {
		(*z) = make(Header, zlqf)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zlqf > 0 {
		var zcua string
		var zxhx interface{}
		zlqf--
		zcua, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		zxhx, bts, err = msgp.ReadIntfBytes(bts)
		if err != nil {
			return
		}
		(*z)[zcua] = zxhx
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Header) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zdaf, zpks := range z {
			_ = zpks
			s += msgp.StringPrefixSize + len(zdaf) + msgp.GuessSize(zpks)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *message) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zeff uint32
	zeff, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zeff > 0 {
		zeff--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "h":
			var zrsw uint32
			zrsw, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Header == nil && zrsw > 0 {
				z.Header = make(Header, zrsw)
			} else if len(z.Header) > 0 {
				for key, _ := range z.Header {
					delete(z.Header, key)
				}
			}
			for zrsw > 0 {
				zrsw--
				var zjfb string
				var zcxo interface{}
				zjfb, err = dc.ReadString()
				if err != nil {
					return
				}
				zcxo, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Header[zjfb] = zcxo
			}
		case "b":
			z.Body, err = dc.ReadBytes(z.Body)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *message) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "h"
	err = en.Append(0x82, 0xa1, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Header)))
	if err != nil {
		return
	}
	for zjfb, zcxo := range z.Header {
		err = en.WriteString(zjfb)
		if err != nil {
			return
		}
		err = en.WriteIntf(zcxo)
		if err != nil {
			return
		}
	}
	// write "b"
	err = en.Append(0xa1, 0x62)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Body)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *message) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "h"
	o = append(o, 0x82, 0xa1, 0x68)
	o = msgp.AppendMapHeader(o, uint32(len(z.Header)))
	for zjfb, zcxo := range z.Header {
		o = msgp.AppendString(o, zjfb)
		o, err = msgp.AppendIntf(o, zcxo)
		if err != nil {
			return
		}
	}
	// string "b"
	o = append(o, 0xa1, 0x62)
	o = msgp.AppendBytes(o, z.Body)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *message) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zxpk uint32
	zxpk, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zxpk > 0 {
		zxpk--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "h":
			var zdnj uint32
			zdnj, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Header == nil && zdnj > 0 {
				z.Header = make(Header, zdnj)
			} else if len(z.Header) > 0 {
				for key, _ := range z.Header {
					delete(z.Header, key)
				}
			}
			for zdnj > 0 {
				var zjfb string
				var zcxo interface{}
				zdnj--
				zjfb, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				zcxo, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Header[zjfb] = zcxo
			}
		case "b":
			z.Body, bts, err = msgp.ReadBytesBytes(bts, z.Body)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *message) Msgsize() (s int) {
	s = 1 + 2 + msgp.MapHeaderSize
	if z.Header != nil {
		for zjfb, zcxo := range z.Header {
			_ = zcxo
			s += msgp.StringPrefixSize + len(zjfb) + msgp.GuessSize(zcxo)
		}
	}
	s += 2 + msgp.BytesPrefixSize + len(z.Body)
	return
}

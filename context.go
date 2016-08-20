package micro

import (
	"time"

	"github.com/nats-io/nats"
	"github.com/plimble/micro/errors"
)

type Context struct {
	*nats.Msg
	Encoder
	pos int8
	mw  []Handler
	c   INats
}

func (c *Context) Next() error {
	//set position to the next
	c.pos++
	midLen := int8(len(c.mw))
	//run the next
	if c.pos < midLen {
		return c.mw[c.pos](c)
	}

	return nil
}

func (c *Context) Publish(subject string, v interface{}) error {
	b, err := c.Encode(v)
	if err != nil {
		return err
	}
	c.c.Publish(subject, b)
	return nil
}

func (c *Context) Request(subject string, req interface{}, res interface{}, timeout time.Duration) error {
	b, err := c.Encode(req)
	if err != nil {
		return err
	}

	msg, err := c.c.Request(subject, b, timeout)
	if err != nil {
		return err
	}

	errProto := &errors.Errors{}
	if err := c.Decode(msg.Data, errProto); err == nil {
		return errProto
	}

	return c.Decode(msg.Data, res)
}

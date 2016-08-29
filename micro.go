package micro

import (
	"os"
	"os/signal"
	"sync"

	"github.com/nats-io/nats"
	perrors "github.com/plimble/errors"
	"github.com/plimble/micro/errors"
)

type Handler func(ctx *Context) error
type ErrorHandler func(ctx *Context, err error) error
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}

//go:generate mockery -name Client -case underscore
type Client interface {
	Publish(subject string, v interface{}, opts ...ClientOption) error
	Request(subject string, req interface{}, res interface{}, opts ...ClientOption) error
	Close()
}

type subHandler struct {
	subject  string
	handlers []Handler
}

type queueSubHandler struct {
	subject  string
	group    string
	handlers []Handler
}

type Micro struct {
	ctxPool sync.Pool
	c       INats
	mw      []Handler
	enc     Encoder
	sub     map[string]subHandler
	qsub    map[string]queueSubHandler
	herr    ErrorHandler
}

func New(c *nats.Conn, enc Encoder) *Micro {
	return &Micro{
		c:    c,
		mw:   []Handler{},
		enc:  enc,
		sub:  make(map[string]subHandler),
		qsub: make(map[string]queueSubHandler),
	}
}

func (m *Micro) Use(h Handler) {
	m.mw = append(m.mw, h)
}

func (m *Micro) HandleError(h ErrorHandler) {
	m.herr = h
}

func (m *Micro) Subscribe(subject string, hs ...Handler) {
	m.sub[subject] = subHandler{
		subject:  subject,
		handlers: joinMiddleware(m.mw, hs),
	}
}

func (m *Micro) QueueSubscribe(subject, group string, hs ...Handler) {
	m.qsub[subject] = queueSubHandler{
		subject:  subject,
		group:    group,
		handlers: joinMiddleware(m.mw, hs),
	}
}

func (m *Micro) RegisterSubscribe() {
	for _, h := range m.sub {
		m.c.Subscribe(h.subject, func(msg *nats.Msg) {
			ms := &message{}
			ms.UnmarshalMsg(msg.Data)

			ctx := m.acquireCtx(msg, m.sub[msg.Subject].handlers, ms.Header, ms.Body, msg.Subject, msg.Reply)
			if err := ctx.Next(); err != nil {
				m.onError(ctx, err)
			}
			m.ctxPool.Put(ctx)
		})
	}
}

func (m *Micro) RegisterQueueSubscribe() {
	for _, h := range m.qsub {
		m.c.QueueSubscribe(h.subject, h.group, func(msg *nats.Msg) {
			ms := &message{}
			ms.UnmarshalMsg(msg.Data)

			ctx := m.acquireCtx(msg, m.qsub[msg.Subject].handlers, ms.Header, ms.Body, msg.Subject, msg.Reply)
			if err := ctx.Next(); err != nil {
				m.onError(ctx, err)
			}
			m.ctxPool.Put(ctx)
		})
	}
}

type HttpError interface {
	Code() int
	Error() string
}

type ProtoError interface {
	StatusCode() int32
	ProtoMessage()
}

func (m *Micro) onError(ctx *Context, err error) {
	if ctx.Reply == "" {
		return
	}

	if m.herr != nil {
		err = m.herr(ctx, err)
	}

	switch werr := perrors.Cause(err).(type) {
	case ProtoError:
		m.Publish(ctx.Reply, werr, nil)
	case HttpError:
		m.Publish(ctx.Reply, errors.New(int32(werr.Code()), werr.Error()), nil)
	default:
		m.Publish(ctx.Reply, errors.New(500, werr.Error()), nil)
	}
}

func (m *Micro) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	m.RegisterSubscribe()
	m.RegisterQueueSubscribe()
	<-c
	m.c.Close()
}

func (m *Micro) acquireCtx(msg *nats.Msg, hs []Handler, header Header, data []byte, subj, reply string) *Context {
	v := m.ctxPool.Get()
	var ctx *Context
	if v == nil {
		ctx = &Context{
			Header:  header,
			Data:    data,
			Subject: subj,
			Reply:   reply,
			Encoder: m.enc,
			mw:      hs,
			Client:  m,
			pos:     -1,
		}
	} else {
		ctx = v.(*Context)
		ctx.Header = header
		ctx.Data = data
		ctx.Subject = subj
		ctx.Reply = reply
		ctx.Encoder = m.enc
		ctx.mw = hs
		ctx.Client = m
		ctx.pos = -1
	}

	return ctx
}

func joinMiddleware(middleware1 []Handler, middleware2 []Handler) []Handler {
	nowLen := len(middleware1)
	totalLen := nowLen + len(middleware2)
	// create a new slice of middleware in order to store all handlers, the already handlers(middleware) and the new
	newMiddleware := make([]Handler, totalLen)
	//copy the already middleware to the just created
	copy(newMiddleware, middleware1)
	//start from there we finish, and store the new middleware too
	copy(newMiddleware[nowLen:], middleware2)
	return newMiddleware
}

func (m *Micro) Publish(subject string, v interface{}, opts ...ClientOption) error {
	o := newOption()
	o.setOptions(opts)

	b, err := m.enc.Encode(v)
	if err != nil {
		return err
	}

	ms := &message{o.header, b}
	mb, err := ms.MarshalMsg(nil)
	if err != nil {
		return err
	}

	m.c.Publish(subject, mb)
	return nil
}

func (m *Micro) Request(subject string, req interface{}, res interface{}, opts ...ClientOption) error {
	o := newOption()
	o.setOptions(opts)

	b, err := m.enc.Encode(req)
	if err != nil {
		return err
	}

	ms := &message{o.header, b}
	mb, err := ms.MarshalMsg(nil)
	if err != nil {
		return err
	}

	msg, err := m.c.Request(subject, mb, o.timeout)
	if err != nil {
		return err
	}

	ms.UnmarshalMsg(msg.Data)

	errProto := &errors.Errors{}
	if err := m.enc.Decode(ms.Body, errProto); err == nil {
		if errProto.Error() != "" {
			return errProto
		}
	}

	return m.enc.Decode(ms.Body, res)
}

func (m *Micro) Forward(subject string, req []byte, opts ...ClientOption) ([]byte, error) {
	o := newOption()
	o.setOptions(opts)

	ms := &message{o.header, req}
	mb, err := ms.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}

	msg, err := m.c.Request(subject, mb, o.timeout)
	if err != nil {
		return nil, err
	}

	ms.UnmarshalMsg(msg.Data)

	errProto := &errors.Errors{}
	if err := m.enc.Decode(ms.Body, errProto); err == nil {
		if errProto.Error() != "" {
			return nil, errProto
		}
	}

	return ms.Body, nil
}

func (m *Micro) Close() {
	m.c.Close()
}

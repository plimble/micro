package micro

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nats-io/nats"
	perrors "github.com/plimble/errors"
	"github.com/plimble/micro/errors"
)

var (
	DefaultTimeout = 2 * time.Second
)

type Handler func(ctx *Context) error
type ErrorHandler func(ctx *Context, err error) error
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}

//go:generate mockery -name Client -case underscore
type Client interface {
	Publish(subject string, v interface{}) error
	Request(subject string, req interface{}, res interface{}, timeout time.Duration) error
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
			ctx := m.acquireCtx(msg, m.sub[msg.Subject].handlers)
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
			ctx := m.acquireCtx(msg, m.qsub[msg.Subject].handlers)
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
		m.Publish(ctx.Reply, werr)
	case HttpError:
		m.Publish(ctx.Reply, errors.New(int32(werr.Code()), werr.Error()))
	default:
		m.Publish(ctx.Reply, errors.New(500, werr.Error()))
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

func (m *Micro) acquireCtx(msg *nats.Msg, hs []Handler) *Context {
	v := m.ctxPool.Get()
	var ctx *Context
	if v == nil {
		ctx = &Context{
			Msg:     msg,
			Encoder: m.enc,
			mw:      hs,
			Client:  m,
			pos:     -1,
		}
	} else {
		ctx = v.(*Context)
		ctx.Msg = msg
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

func (m *Micro) Publish(subject string, v interface{}) error {
	b, err := m.enc.Encode(v)
	if err != nil {
		return err
	}
	m.c.Publish(subject, b)
	return nil
}

func (m *Micro) Request(subject string, req interface{}, res interface{}, timeout time.Duration) error {
	b, err := m.enc.Encode(req)
	if err != nil {
		return err
	}

	msg, err := m.c.Request(subject, b, timeout)
	if err != nil {
		return err
	}

	errProto := &errors.Errors{}
	if err := m.enc.Decode(msg.Data, errProto); err == nil {
		if errProto.Error() != "" {
			return errProto
		}
	}

	return m.enc.Decode(msg.Data, res)
}

func (m *Micro) Close() {
	m.c.Close()
}

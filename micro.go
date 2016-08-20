package micro

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats"
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

type Micro struct {
	ctxPool sync.Pool
	c       INats
	mw      []Handler
	enc     Encoder
	sub     map[string][]Handler
	qsub    map[string][]Handler
	herr    ErrorHandler
}

func New(c *nats.Conn, enc Encoder) *Micro {
	return &Micro{
		c:    c,
		mw:   []Handler{},
		enc:  enc,
		sub:  make(map[string][]Handler),
		qsub: make(map[string][]Handler),
	}
}

func (m *Micro) Use(h Handler) {
	m.mw = append(m.mw, h)
}

func (m *Micro) HandleError(h ErrorHandler) {
	m.herr = h
}

func (m *Micro) Subscribe(subject string, hs ...Handler) {
	newhs := joinMiddleware(m.mw, hs)
	m.sub[subject] = newhs
}

func (m *Micro) QueueSubscribe(subject, group string, hs ...Handler) {
	newhs := joinMiddleware(m.mw, hs)
	m.qsub[subject+"|"+group] = newhs
}

func (m *Micro) RegisterSubscribe() {
	for subj, hs := range m.sub {
		m.c.Subscribe(subj, func(msg *nats.Msg) {
			ctx := m.acquireCtx(msg, hs)
			if err := ctx.Next(); err != nil {
				m.onError(ctx, err)
			}
			m.ctxPool.Put(ctx)
		})
	}
}

func (m *Micro) RegisterQueueSubscribe() {
	for qsubj, hs := range m.qsub {
		subj := strings.Split(qsubj, "|")
		m.c.QueueSubscribe(subj[0], subj[1], func(msg *nats.Msg) {
			ctx := m.acquireCtx(msg, hs)
			if err := ctx.Next(); err != nil {
				m.onError(ctx, err)
			}
			m.ctxPool.Put(ctx)
		})
	}
}

type HttpProtoError interface {
	StatusCode() int32
	Error() string
	ProtoMessage()
}

func (m *Micro) onError(ctx *Context, err error) {
	if ctx.Reply == "" {
		return
	}

	if m.herr != nil {
		err = m.herr(ctx, err)
	}

	switch werr := err.(type) {
	case HttpProtoError:
		m.Publish(ctx.Reply, werr)
	default:
		m.Publish(ctx.Reply, errors.New(500, err.Error()))
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
			c:       m.c,
			pos:     -1,
		}
	} else {
		ctx = v.(*Context)
		ctx.mw = hs
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
		return errProto
	}

	return m.enc.Decode(msg.Data, res)
}

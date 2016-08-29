package micro

import "time"

var (
	DefaultTimeout = 2 * time.Second
)

type ClientOption func(o *option)

type option struct {
	timeout time.Duration
	header  Header
}

func newOption() *option {
	return &option{
		timeout: DefaultTimeout,
	}
}

func (o *option) setOptions(opts []ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *option) {
		o.timeout = timeout
	}
}

func WithHeader(header Header) ClientOption {
	return func(o *option) {
		o.header = header
	}
}

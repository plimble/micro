package micro

import (
	"time"

	"github.com/nats-io/nats"
)

type INats interface {
	Subscribe(subj string, msg nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribe(subj, gtroup string, msg nats.MsgHandler) (*nats.Subscription, error)
	Publish(subj string, data []byte) error
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	Close()
}

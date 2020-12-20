package messages

import "context"

type Message interface {
	Ack()
	Nack()
	Data() []byte
}

type ListenerFuncType func(context.Context, Message) error

type Listener interface {
	Listen(ctx context.Context, topic string, callback ListenerFuncType) error
}

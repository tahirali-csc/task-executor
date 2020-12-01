package pubsub

import "context"

type Message interface{}

type Pubsub interface {
	NewTopic(topic string)
	Publish(ctx context.Context, topic string, msg Message) error
	Subscribe(ctxt context.Context, topic string) (chan Message, error)
}

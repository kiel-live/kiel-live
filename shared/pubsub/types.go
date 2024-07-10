package pubsub

import (
	"context"
)

// Message defines a published message.
type Message []byte

// Subscriber receives published messages.
type Subscriber func(Message)

type Broker interface {
	Publish(c context.Context, topic string, message Message) error
	Subscribe(c context.Context, topic string, subscriber Subscriber) error
}

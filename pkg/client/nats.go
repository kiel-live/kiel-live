package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/protocol"
	"github.com/nats-io/nats.go"
)

type natsClient struct {
	nc                       *nats.Conn
	JS                       nats.JetStreamContext
	subscriptions            map[string]*nats.Subscription
	host                     string
	username                 string
	password                 string
	connectionHandler        func(connected bool)
	topicSubscriptionHandler func(topic []string)

	topicSubscriptions map[string][]string
	subscriptionsMu    sync.Mutex
}

type NatsOption func(c *natsClient)

func NewNatsClient(host string, opts ...NatsOption) Client {
	client := &natsClient{
		subscriptions: make(map[string]*nats.Subscription),
		host:          host,
		username:      "",
		password:      "",
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithAuth(username string, password string) NatsOption {
	return func(c *natsClient) {
		c.username = username
		c.password = password
	}
}

func (n *natsClient) Connect() (err error) {
	if len(n.username) < 1 && len(n.password) < 1 {
		n.nc, err = nats.Connect(n.host)
	} else {
		n.nc, err = nats.Connect(n.host, nats.UserInfo(n.username, n.password))
	}

	if err != nil {
		return err
	}

	n.JS, err = n.nc.JetStream()
	return err
}

func (n *natsClient) IsConnected() bool {
	return n.nc.IsConnected()
}

// Close will unsubscribe all subjects and shutdown connection
func (n *natsClient) Disconnect() error {
	for subject := range n.subscriptions {
		err := n.Unsubscribe(subject)
		if err != nil {
			return err
		}
	}
	n.nc.Close()
	return nil
}

func (n *natsClient) Subscribe(subject string, cb SubscribeCallback) error {
	if n.subscriptions[subject] != nil {
		return fmt.Errorf("already subscribed to '%s'", subject)
	}

	sub, err := n.nc.Subscribe(subject, func(msg *nats.Msg) {
		cb(&TopicMessage{
			Topic: msg.Subject,
			Data:  string(msg.Data),
			//Raw:     msg,
		})
	})
	if err != nil {
		return err
	}

	n.subscriptions[subject] = sub

	return nil
}

func (n *natsClient) Unsubscribe(subject string) error {
	sub := n.subscriptions[subject]
	if sub != nil {
		return fmt.Errorf("you have not subscribed to that subject '%s'", subject)
	}

	msg, err := n.nc.Request(protocol.TopicRequestUnsubscribe, []byte(subject), 1*time.Second)
	if err != nil {
		return err
	}

	if string(msg.Data) != "ok" {
		return fmt.Errorf("unsubscribe failed '%s'", subject)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	delete(n.subscriptions, subject)

	return nil
}

func (n *natsClient) Publish(subject string, data string) error {
	return n.nc.Publish(subject, []byte(data))
}

func (n *natsClient) SetConnectionHandler(connectionHandler func(connected bool)) {
	n.connectionHandler = connectionHandler
}

func (n *natsClient) SetTopicSubscriptionHandler(topicSubscriptionHandler func(topics []string)) {
	n.topicSubscriptionHandler = topicSubscriptionHandler
}

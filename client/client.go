package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/kiel-live/kiel-live/protocol"
	"github.com/nats-io/nats.go"
)

// WebSocketClient return websocket client connection
type Client struct {
	nc            *nats.Conn
	JS            nats.JetStreamContext
	subscriptions map[string]*nats.Subscription
	host          string
	username      string
	password      string
}

type Option func(c *Client)

// NewClient create new connection
func NewClient(host string, opts ...Option) *Client {
	client := &Client{
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

func WithAuth(username string, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func (c *Client) Connect() (err error) {
	if len(c.username) < 1 && len(c.password) < 1 {
		c.nc, err = nats.Connect(c.host)
	} else {
		c.nc, err = nats.Connect(c.host, nats.UserInfo(c.username, c.password))
	}

	if err != nil {
		return err
	}

	c.JS, err = c.nc.JetStream()
	return err
}

func (c *Client) IsConnected() bool {
	return c.nc.IsConnected()
}

// Close will unsubscribe all topics and shutdown connection
func (c *Client) Disconnect() error {
	for topic := range c.subscriptions {
		err := c.Unsubscribe(topic)
		if err != nil {
			return err
		}
	}
	c.nc.Close()
	return nil
}

type TopicMessage struct {
	Topic string
	Data  string
	Raw   *nats.Msg
}
type SubscribeCallback func(msg *TopicMessage)
type SubscribeOption func(topic string, cb SubscribeCallback) error

func (c *Client) Subscribe(topic string, cb SubscribeCallback, opts ...SubscribeOption) error {
	if c.subscriptions[topic] != nil {
		return fmt.Errorf("Already subscribed to '%s'", topic)
	}

	sub, err := c.nc.Subscribe(topic, func(msg *nats.Msg) {
		cb(&TopicMessage{
			Topic: msg.Subject,
			Data:  string(msg.Data),
			Raw:   msg,
		})
	})
	if err != nil {
		return err
	}

	c.subscriptions[topic] = sub

	for _, opt := range opts {
		err := opt(topic, cb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) WithAck() SubscribeOption {
	return func(topic string, _ SubscribeCallback) error {
		msg, err := c.nc.Request(protocol.TopicRequestSubscribe, []byte(topic), 1*time.Second)
		if err != nil {
			if err.Error() == "nats: timeout" {
				return fmt.Errorf("No one is answering us")
			}
			return err
		}

		if !strings.HasPrefix(string(msg.Data), "ok") {
			return fmt.Errorf("Can't subscribe to '%s'", topic)
		}

		return nil
	}
}

func (c *Client) WithCache() SubscribeOption {
	return func(topic string, cb SubscribeCallback) error {
		msg, err := c.nc.Request(protocol.TopicRequestCache, []byte(topic), 1*time.Second)
		if err != nil {
			// TODO ignore cache miss or timeouts
			// return err
			return nil
		}

		data := string(msg.Data)
		if data == "err" {
			// return fmt.Errorf("Miss on cache")
			return nil
		}

		cb(&TopicMessage{
			Topic: msg.Subject,
			Data:  data,
			Raw:   msg,
		})

		return nil
	}
}

func (c *Client) Unsubscribe(topic string) error {
	sub := c.subscriptions[topic]
	if sub != nil {
		return fmt.Errorf("You have not subscribed to that topic '%s'", topic)
	}

	msg, err := c.nc.Request(protocol.TopicRequestUnsubscribe, []byte(topic), 1*time.Second)
	if err != nil {
		return err
	}

	if string(msg.Data) != "ok" {
		return fmt.Errorf("Unsubscription failed '%s'", topic)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	delete(c.subscriptions, topic)

	return nil
}

func (c *Client) Publish(topic string, data string) error {
	return c.PublishRaw(topic, []byte(data))
}

func (c *Client) PublishRaw(topic string, data []byte) error {
	return c.nc.Publish(topic, data)
}

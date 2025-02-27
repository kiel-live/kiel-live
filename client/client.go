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

// Close will unsubscribe all subjects and shutdown connection
func (c *Client) Disconnect() error {
	for subject := range c.subscriptions {
		err := c.Unsubscribe(subject)
		if err != nil {
			return err
		}
	}
	c.nc.Close()
	return nil
}

type SubjectMessage struct {
	Subject string
	Data    string
	Raw     *nats.Msg
}
type SubscribeCallback func(msg *SubjectMessage)
type SubscribeOption func(subject string, cb SubscribeCallback) error

func (c *Client) Subscribe(subject string, cb SubscribeCallback, opts ...SubscribeOption) error {
	if c.subscriptions[subject] != nil {
		return fmt.Errorf("Already subscribed to '%s'", subject)
	}

	sub, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		cb(&SubjectMessage{
			Subject: msg.Subject,
			Data:    string(msg.Data),
			Raw:     msg,
		})
	})
	if err != nil {
		return err
	}

	c.subscriptions[subject] = sub

	for _, opt := range opts {
		err := opt(subject, cb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) WithAck() SubscribeOption {
	return func(subject string, _ SubscribeCallback) error {
		msg, err := c.nc.Request(protocol.SubjectRequestSubscribe, []byte(subject), 1*time.Second)
		if err != nil {
			if err.Error() == "nats: timeout" {
				return fmt.Errorf("No one is answering us")
			}
			return err
		}

		if !strings.HasPrefix(string(msg.Data), "ok") {
			return fmt.Errorf("Can't subscribe to '%s'", subject)
		}

		return nil
	}
}

func (c *Client) WithCache() SubscribeOption {
	return func(subject string, cb SubscribeCallback) error {
		msg, err := c.nc.Request(protocol.SubjectRequestCache, []byte(subject), 1*time.Second)
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

		cb(&SubjectMessage{
			Subject: msg.Subject,
			Data:    data,
			Raw:     msg,
		})

		return nil
	}
}

func (c *Client) Unsubscribe(subject string) error {
	sub := c.subscriptions[subject]
	if sub != nil {
		return fmt.Errorf("You have not subscribed to that subject '%s'", subject)
	}

	msg, err := c.nc.Request(protocol.SubjectRequestUnsubscribe, []byte(subject), 1*time.Second)
	if err != nil {
		return err
	}

	if string(msg.Data) != "ok" {
		return fmt.Errorf("Unsubscription failed '%s'", subject)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	delete(c.subscriptions, subject)

	return nil
}

func (c *Client) Publish(subject string, data string) error {
	return c.PublishRaw(subject, []byte(data))
}

func (c *Client) PublishRaw(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

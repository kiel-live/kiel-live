package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/pkg/models"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type natsClient struct {
	nc                       *nats.Conn
	JS                       nats.JetStreamContext
	subscriptions            map[string]*nats.Subscription // active subscriptions by this client
	host                     string
	username                 string
	password                 string
	connectionHandler        func(connected bool)
	topicSubscriptionHandler func(topic string, added bool)

	topicSubscriptions map[string][]string // topics on the server subscribed to by some client
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
	opts := []nats.Option{
		nats.Name("Kiel Live Collector"),
		nats.ConnectHandler(func(conn *nats.Conn) {
			if conn.IsConnected() {
				n.initTopics()
				logrus.Debugf("Connected to NATS server at %s", n.host)
			}
			if n.connectionHandler != nil {
				n.connectionHandler(conn.IsConnected())
			}
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			if n.connectionHandler != nil {
				n.connectionHandler(conn.IsConnected())
			}
		}),
	}

	if n.username != "" && n.password != "" {
		opts = append(opts, nats.UserInfo(n.username, n.password))
	}

	n.nc, err = nats.Connect(n.host, opts...)

	if err != nil {
		return err
	}

	n.JS, err = n.nc.JetStream()
	return err
}

func (n *natsClient) IsConnected() bool {
	return n.nc.IsConnected()
}

// Close will unsubscribe all topics and shutdown connection
func (n *natsClient) Disconnect() error {
	for topic := range n.subscriptions {
		err := n.Unsubscribe(topic)
		if err != nil {
			return err
		}
	}
	n.nc.Close()
	return nil
}

func (n *natsClient) Subscribe(topic string, cb SubscribeCallback) error {
	if _, ok := n.subscriptions[topic]; ok {
		return fmt.Errorf("already subscribed to '%s'", topic)
	}

	sub, err := n.nc.Subscribe(topic, func(msg *nats.Msg) {
		cb(&Message{
			Topic: msg.Subject,
			Data:  string(msg.Data),
		})
	})
	if err != nil {
		return err
	}

	n.subscriptions[topic] = sub

	return nil
}

func (n *natsClient) Unsubscribe(topic string) error {
	sub := n.subscriptions[topic]
	if sub != nil {
		return fmt.Errorf("you have not subscribed to that topic '%s'", topic)
	}

	msg, err := n.nc.Request(models.TopicRequestUnsubscribe, []byte(topic), 1*time.Second)
	if err != nil {
		return err
	}

	if string(msg.Data) != "ok" {
		return fmt.Errorf("unsubscribe failed '%s'", topic)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return err
	}

	delete(n.subscriptions, topic)

	return nil
}

func (n *natsClient) SetOnConnectionChanged(connectionHandler func(connected bool)) {
	n.connectionHandler = connectionHandler
}

func (n *natsClient) SetOnTopicsChanged(topicSubscriptionHandler func(topic string, added bool)) {
	n.topicSubscriptionHandler = topicSubscriptionHandler
}

func (n *natsClient) UpdateStop(stop *models.Stop) error {
	jsonData, err := json.Marshal(stop)
	if err != nil {
		return err
	}

	return n.nc.Publish(fmt.Sprintf(models.TopicStop, stop.ID), jsonData)
}

func (n *natsClient) UpdateVehicle(vehicle *models.Vehicle) error {
	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		return err
	}

	return n.nc.Publish(fmt.Sprintf(models.TopicVehicle, vehicle.ID), jsonData)
}

func (n *natsClient) UpdateTrip(trip *models.Trip) error {
	jsonData, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	return n.nc.Publish(fmt.Sprintf(models.TopicTrip, trip.ID), jsonData)
}

func (n *natsClient) DeleteStop(stopID string) error {
	return n.nc.Publish(fmt.Sprintf(models.TopicStop, stopID), []byte(models.DeletePayload))
}

func (n *natsClient) DeleteVehicle(vehicleID string) error {
	return n.nc.Publish(fmt.Sprintf(models.TopicVehicle, vehicleID), []byte(models.DeletePayload))
}

func (n *natsClient) DeleteTrip(tripID string) error {
	return n.nc.Publish(fmt.Sprintf(models.TopicTrip, tripID), []byte(models.DeletePayload))
}

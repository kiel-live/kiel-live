package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/pkg/models"

	"github.com/nats-io/nats.go"
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

	closed    chan struct{} // closed once the connection is gone for good
	closeOnce sync.Once

	advisoryOnce sync.Once // guards the one-time consumer advisory subscriptions
}

type NatsOption func(c *natsClient)

func NewNatsClient(host string, opts ...NatsOption) Client {
	client := &natsClient{
		subscriptions: make(map[string]*nats.Subscription),
		host:          host,
		username:      "",
		password:      "",
		closed:        make(chan struct{}),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func NatsWithAuth(username string, password string) NatsOption {
	return func(c *natsClient) {
		c.username = username
		c.password = password
	}
}

func (n *natsClient) Connect() (err error) {
	onConnected := func(conn *nats.Conn) {
		if conn.IsConnected() {
			n.initTopics()
		}
		if n.connectionHandler != nil {
			n.connectionHandler(conn.IsConnected())
		}
	}

	opts := []nats.Option{
		nats.Name("Kiel Live Collector"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.ReconnectJitter(500*time.Millisecond, time.Second),
		nats.ConnectHandler(func(conn *nats.Conn) {
			slog.Debug("Connected to NATS server", "host", n.host)
			onConnected(conn)
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			slog.Warn("Disconnected from NATS server, reconnecting ...", "host", n.host, "error", err)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			slog.Info("Reconnected to NATS server", "host", n.host)
			onConnected(conn)
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			n.closeOnce.Do(func() {
				close(n.closed)
			})
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

func (n *natsClient) Closed() <-chan struct{} {
	return n.closed
}

// Close will unsubscribe all topics and shutdown connection
func (n *natsClient) Disconnect() error {
	if n.nc.IsClosed() {
		return nil
	}

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
		data := json.RawMessage(msg.Data)
		cb(&Message{
			Topic:  msg.Subject,
			Action: "",
			Data:   &data,
			SentAt: time.Now(),
		})
	})
	if err != nil {
		return err
	}

	n.subscriptions[topic] = sub

	return nil
}

func (n *natsClient) Unsubscribe(topic string) error {
	sub, ok := n.subscriptions[topic]
	if !ok || sub == nil {
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

// parseOptionalTime parses an optional RFC3339 timestamp. An empty string yields the zero time without an error.
func parseOptionalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, value)
}

func (n *natsClient) UpdateStop(stop *models.Stop) error {
	// TODO: remove once majority of clients updated (added 29.03.2026)
	stop.Arrivals = nil
	for _, departure := range stop.Departures {
		log := slog.With(
			"stop_id", stop.ID,
			"trip_id", departure.TripID,
			"route_name", departure.RouteName,
			"name", departure.Name,
		)

		planned, err := parseOptionalTime(departure.Planned)
		if err != nil {
			log.Warn("Stop departure has an invalid planned time, skipping", "planned", departure.Planned, "error", err)
			continue
		}

		actual, err := parseOptionalTime(departure.Actual)
		if err != nil {
			log.Warn("Stop departure has an invalid actual time, skipping", "actual", departure.Actual, "error", err)
			continue
		}

		// not every departure comes with a planned time, fall back to the actual one
		if planned.IsZero() {
			planned = actual
		}

		if planned.IsZero() {
			log.Warn("Stop departure has neither a planned nor an actual time, skipping")
			continue
		}

		// legacy clients expect the seconds until arrival, 0 means "unknown"
		eta := 0
		if !actual.IsZero() {
			eta = int(time.Until(actual).Seconds())
		}

		arrival := &models.StopArrival{ //nolint:staticcheck
			Name:      departure.Name,
			Type:      departure.Type,
			VehicleID: departure.VehicleID,
			TripID:    departure.TripID,
			RouteID:   departure.RouteID,
			RouteName: departure.RouteName,
			Direction: departure.Direction,
			State:     string(departure.State),
			Platform:  departure.Platform,
			Planned:   planned.Format("15:04"),
			Eta:       eta,
		}
		stop.Arrivals = append(stop.Arrivals, arrival)
	}

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
	// TODO: remove once majority of clients updated (added 29.03.2026)
	trip.Arrivals = nil
	for _, departure := range trip.Departures {
		log := slog.With("trip_id", trip.ID, "name", departure.Name)

		planned, err := parseOptionalTime(departure.Planned)
		if err != nil {
			log.Warn("Trip departure has an invalid planned time, skipping", "planned", departure.Planned, "error", err)
			continue
		}

		// not every departure comes with a planned time, fall back to the actual one
		if planned.IsZero() {
			actual, err := parseOptionalTime(departure.Actual)
			if err != nil {
				log.Warn("Trip departure has an invalid actual time, skipping", "actual", departure.Actual, "error", err)
				continue
			}

			planned = actual
		}

		if planned.IsZero() {
			log.Warn("Trip departure has neither a planned nor an actual time, skipping")
			continue
		}

		arrival := &models.TripArrival{ //nolint:staticcheck
			Name:    departure.Name,
			State:   string(departure.State),
			Planned: planned.Format("15:04"),
		}
		trip.Arrivals = append(trip.Arrivals, arrival)
	}

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

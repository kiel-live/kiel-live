package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/protocol"
)

const (
	hubReconnectDelay = 5 * time.Second
)

type hubClient struct {
	url   string
	token string

	writeMu sync.Mutex
	conn    *websocket.Conn

	connected int32 // atomic bool

	topicsMu         sync.Mutex
	subscribedTopics map[string]bool // NATS-format topics: "data.map.stop.kvg-123"

	connectionHandler        func(bool)
	topicSubscriptionHandler func(string, bool)

	stopCh chan struct{}
}

func NewHubClient(url, token string) Client {
	return &hubClient{
		url:              url,
		token:            token,
		subscribedTopics: make(map[string]bool),
		stopCh:           make(chan struct{}),
	}
}

func (c *hubClient) dial() (*websocket.Conn, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)

	header := http.Header{"Authorization": {"Bearer " + c.token}}
	conn, _, err := websocket.DefaultDialer.Dial(u.String()+"/ws/collector", header)
	return conn, err
}

func (c *hubClient) onConnect(conn *websocket.Conn) error {
	env := protocol.Envelope{
		Method: protocol.MethodSubscribe,
		Topic:  protocol.TopicSubscriptions,
	}
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}

func (c *hubClient) Connect() error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	if err := c.onConnect(conn); err != nil {
		conn.Close()
		return err
	}
	c.writeMu.Lock()
	c.conn = conn
	c.writeMu.Unlock()
	atomic.StoreInt32(&c.connected, 1)

	go c.readLoop(conn)
	return nil
}

func (c *hubClient) readLoop(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if atomic.LoadInt32(&c.connected) == 0 {
				return // intentional disconnect
			}
			slog.Warn("hub client connection lost, reconnecting", "error", err)
			conn.Close()
			atomic.StoreInt32(&c.connected, 0)
			if c.connectionHandler != nil {
				c.connectionHandler(false)
			}

			// clear topics on disconnect
			c.topicsMu.Lock()
			c.subscribedTopics = make(map[string]bool)
			c.topicsMu.Unlock()

			// reconnect loop
			for {
				select {
				case <-c.stopCh:
					return
				case <-time.After(hubReconnectDelay):
				}
				newConn, err := c.dial()
				if err != nil {
					slog.Warn("hub client reconnect failed", "error", err)
					continue
				}
				if err := c.onConnect(newConn); err != nil {
					newConn.Close()
					slog.Warn("hub client reconnect onConnect failed", "error", err)
					continue
				}
				c.writeMu.Lock()
				c.conn = newConn
				c.writeMu.Unlock()
				conn = newConn
				atomic.StoreInt32(&c.connected, 1)
				slog.Info("hub client reconnected")
				if c.connectionHandler != nil {
					c.connectionHandler(true)
				}
				break
			}
			continue
		}

		var env protocol.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			slog.Error("hub client parse error", "error", err)
			continue
		}
		c.handleEnvelope(env)
	}
}

func (c *hubClient) handleEnvelope(env protocol.Envelope) {
	if env.Method != protocol.MethodUpdate || env.Topic != protocol.TopicSubscriptions || env.Data == nil {
		return
	}
	var se protocol.SubscriptionEvent
	if err := json.Unmarshal(*env.Data, &se); err != nil {
		slog.Error("hub client subscription event parse error", "error", err)
		return
	}

	natsTopic := protocol.HubTopicToNats(se.Topic)

	c.topicsMu.Lock()
	if se.Subscribed {
		c.subscribedTopics[natsTopic] = true
	} else {
		delete(c.subscribedTopics, natsTopic)
	}
	c.topicsMu.Unlock()

	if c.topicSubscriptionHandler != nil {
		c.topicSubscriptionHandler(natsTopic, se.Subscribed)
	}
}

func (c *hubClient) sendEnvelope(env protocol.Envelope) error {
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.conn == nil {
		return nil // silently drop while disconnected
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *hubClient) IsConnected() bool {
	return atomic.LoadInt32(&c.connected) == 1
}

func (c *hubClient) Disconnect() error {
	atomic.StoreInt32(&c.connected, 0)
	close(c.stopCh)
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *hubClient) SetOnConnectionChanged(handler func(connected bool)) {
	c.connectionHandler = handler
}

func (c *hubClient) GetSubscribedTopics() []string {
	c.topicsMu.Lock()
	defer c.topicsMu.Unlock()
	topics := make([]string, 0, len(c.subscribedTopics))
	for t := range c.subscribedTopics {
		topics = append(topics, t)
	}
	return topics
}

func (c *hubClient) SetOnTopicsChanged(handler func(topic string, subscribed bool)) {
	c.topicSubscriptionHandler = handler
}

func (c *hubClient) UpdateStop(stop *models.Stop) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodUpdate,
		Topic:  protocol.FormatDetailTopic("stop", stop.ID),
		Data:   toRawMessage(stop),
	})
}

func (c *hubClient) UpdateVehicle(vehicle *models.Vehicle) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodUpdate,
		Topic:  protocol.FormatDetailTopic("vehicle", vehicle.ID),
		Data:   toRawMessage(vehicle),
	})
}

func (c *hubClient) UpdateTrip(trip *models.Trip) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodUpdate,
		Topic:  protocol.FormatDetailTopic("trip", trip.ID),
		Data:   toRawMessage(trip),
	})
}

func (c *hubClient) DeleteStop(stopID string) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodDelete,
		Topic:  protocol.FormatDetailTopic("stop", stopID),
	})
}

func (c *hubClient) DeleteVehicle(vehicleID string) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodDelete,
		Topic:  protocol.FormatDetailTopic("vehicle", vehicleID),
	})
}

func (c *hubClient) DeleteTrip(tripID string) error {
	return c.sendEnvelope(protocol.Envelope{
		Method: protocol.MethodDelete,
		Topic:  protocol.FormatDetailTopic("trip", tripID),
	})
}

func toRawMessage(v any) *json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	raw := json.RawMessage(data)
	return &raw
}

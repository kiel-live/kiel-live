package hub

import (
	"errors"
	"sync"

	"github.com/kiel-live/kiel-live/backend/proto"
)

type subscriptionRequest struct {
	Connection connection
	Channel    string
	Done       chan error
}

type Hub struct {
	quit chan struct{}

	// Keeps track of all channels a connection is subscribed to.
	subscriptions map[connection]map[string]bool

	// Allows mapping channels to subscribers.
	channels map[string]map[connection]bool

	// Makes tokens to connections
	connections map[string]connection

	newSubscriptions   chan subscriptionRequest
	newUnsubscriptions chan subscriptionRequest
	newChannelMessage  chan proto.ClientMessage

	sync.Mutex
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) Prepare() error {
	h.quit = make(chan struct{})

	h.subscriptions = make(map[connection]map[string]bool)
	h.channels = make(map[string]map[connection]bool)
	h.connections = make(map[string]connection)

	h.newSubscriptions = make(chan subscriptionRequest, 100)
	h.newUnsubscriptions = make(chan subscriptionRequest, 100)

	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case r := <-h.newSubscriptions:
			h.handleSubscribe(r)
		case r := <-h.newUnsubscriptions:
			h.handleUnsubscribe(r)
		case m := <-h.newChannelMessage:
			h.handleMessage(m)
		case <-h.quit:
			return
		}
	}
}

func (h *Hub) Stop() {
	h.quit <- struct{}{}
}

func (h *Hub) Connect(conn connection) error {
	h.Lock()
	defer h.Unlock()

	h.subscriptions[conn] = make(map[string]bool)
	h.connections[conn.GetToken()] = conn
	return nil
}

func (h *Hub) Disconnect(conn connection) error {
	if !h.hasConnection(conn) {
		return errors.New("Unknown connection")
	}

	h.Lock()
	channels := make([]string, 0)
	for channel, _ := range h.subscriptions[conn] {
		channels = append(channels, channel)
	}
	h.Unlock()

	// Unsubscribe from all channels
	for _, channel := range channels {
		err := h.Unsubscribe(conn, channel)
		if err != nil {
			return err
		}
	}

	h.Lock()
	defer h.Unlock()
	delete(h.subscriptions, conn)
	delete(h.connections, conn.GetToken())
	return nil
}

func (h *Hub) hasConnection(conn connection) bool {
	h.Lock()
	defer h.Unlock()

	_, ok := h.subscriptions[conn]
	return ok
}

func (h *Hub) hasSubscription(conn connection, channel string) bool {
	h.Lock()
	defer h.Unlock()

	s, ok := h.subscriptions[conn]
	if !ok {
		return false
	}

	_, ok = s[channel]
	return ok
}

func (h *Hub) Subscribe(conn connection, channel string) error {
	if !h.hasConnection(conn) {
		return errors.New("Unknown connection")
	}

	r := subscriptionRequest{
		Connection: conn,
		Channel:    channel,
		Done:       make(chan error),
	}
	h.newSubscriptions <- r
	return <-r.Done
}

func (h *Hub) handleSubscribe(r subscriptionRequest) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.channels[r.Channel]; !ok {
		// New channel! Try to connect to Redis first
		h.channels[r.Channel] = make(map[connection]bool)
	}

	h.subscriptions[r.Connection][r.Channel] = true
	h.channels[r.Channel][r.Connection] = true
	r.Done <- nil
}

func (h *Hub) Unsubscribe(conn connection, channel string) error {
	if !h.hasConnection(conn) {
		return errors.New("Unknown connection")
	}
	if !h.hasSubscription(conn, channel) {
		// Some clients seem to be sending double unsubscribes,
		// ignore those for now:
		//return fmt.Errorf("Not subscribed to channel %s", channel)
		return nil
	}

	r := subscriptionRequest{
		Connection: conn,
		Channel:    channel,
		Done:       make(chan error),
	}
	h.newUnsubscriptions <- r
	return <-r.Done
}

func (h *Hub) handleUnsubscribe(r subscriptionRequest) {
	h.Lock()
	defer h.Unlock()

	delete(h.subscriptions[r.Connection], r.Channel)
	delete(h.channels[r.Channel], r.Connection)

	if len(h.channels[r.Channel]) == 0 {
		// Last subscriber, release it.
		delete(h.channels, r.Channel)
	}

	r.Done <- nil
}

// func (h *Hub) processClient(t, token string, args []string) {
// 	if c, ok := h.connections[token]; ok {
// 		c.Process(t, args)
// 	}
// }

func (h *Hub) handleMessage(m proto.ClientMessage) {
	h.Lock()
	defer h.Unlock()
	channel := m.Channel()
	data := m.Data()

	if _, ok := h.channels[channel]; !ok {
		return // No longer subscribed?
	}

	for conn, _ := range h.channels[channel] {
		conn.Send(channel, data)
	}
}

type hubStats struct {
	LocalSubscriptions map[string]int
}

func (h *Hub) Stats() (hubStats, error) {
	h.Lock()
	defer h.Unlock()

	subscriptions := make(map[string]int)
	for k, v := range h.channels {
		subscriptions[k] = len(v)
	}

	return hubStats{
		LocalSubscriptions: subscriptions,
	}, nil
}

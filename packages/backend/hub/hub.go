package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/kiel-live/kiel-live/packages/backend/store"
	protocol "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("origin", "Hub")

type subscriptionRequest struct {
	Connection connection
	Channel    string
	Done       chan error
}

type channelMessageRequest struct {
	Connection connection
	Channel    string
	Data       string
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
	newChannelMessages chan channelMessageRequest

	store store.Store

	sync.Mutex
}

func NewHub(store store.Store) (*Hub, error) {
	hub := &Hub{
		store: store,
	}

	err := hub.Prepare()
	if err != nil {
		return nil, err
	}

	return hub, nil
}

func (h *Hub) Prepare() error {
	h.quit = make(chan struct{})

	h.subscriptions = make(map[connection]map[string]bool)
	h.channels = make(map[string]map[connection]bool)
	h.connections = make(map[string]connection)

	h.newSubscriptions = make(chan subscriptionRequest, 100)
	h.newUnsubscriptions = make(chan subscriptionRequest, 100)
	h.newChannelMessages = make(chan channelMessageRequest, 100)

	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case r := <-h.newSubscriptions:
			h.handleSubscribe(r)
		case r := <-h.newUnsubscriptions:
			h.handleUnsubscribe(r)
		case m := <-h.newChannelMessages:
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
		return errors.New("unknown connection")
	}

	h.Lock()
	channels := make([]string, 0)
	for channel := range h.subscriptions[conn] {
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
		return errors.New("unknown connection")
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
		// New channel!
		h.channels[r.Channel] = make(map[connection]bool)
		h.publishListOfChannels()
	}

	h.subscriptions[r.Connection][r.Channel] = true
	h.channels[r.Channel][r.Connection] = true

	// send cached data of channel to client if it exists
	data, err := h.store.Get(r.Channel)
	if err == nil {
		r.Connection.Send(r.Channel, data)
	}

	r.Done <- nil
}

func (h *Hub) Unsubscribe(conn connection, channel string) error {
	if !h.hasConnection(conn) {
		return errors.New("unknown connection")
	}
	if !h.hasSubscription(conn, channel) {
		return fmt.Errorf("not subscribed to channel %s", channel)
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
		// Last subscriber, release channel.
		delete(h.channels, r.Channel)
		h.publishListOfChannels()
	}

	r.Done <- nil
}

func (h *Hub) publishListOfChannels() {
	channelNames := make([]string, len(h.channels))

	i := 0
	for k := range h.channels {
		channelNames[i] = k
		i++
	}

	data, err := json.Marshal(channelNames)
	if err != nil {
		return
	}

	r := channelMessageRequest{
		Connection: nil,
		Channel:    protocol.ChannelNameSubscribedChannels,
		Data:       string(data),
		Done:       nil,
	}
	h.newChannelMessages <- r
}

func (h *Hub) Publish(conn connection, channel string, data string) error {
	if !h.hasConnection(conn) {
		return errors.New("unknown connection")
	}

	if conn != nil && channel == protocol.ChannelNameSubscribedChannels {
		return errors.New("you are not allowed to publish to this channel")
	}

	r := channelMessageRequest{
		Connection: conn,
		Channel:    channel,
		Data:       data,
		Done:       make(chan error),
	}
	h.newChannelMessages <- r
	return <-r.Done
}

func (h *Hub) handleMessage(m channelMessageRequest) {
	h.Lock()
	defer h.Unlock()
	channel := m.Channel
	data := m.Data

	if storedData, _ := h.store.Get(channel); storedData == data {
		log.Infoln("Data of %s has already been published. Skipping ...")
		if m.Done != nil {
			m.Done <- nil
		}
		return
	}

	// write data to cache of data-store
	err := h.store.Set(channel, data)
	if err != nil {
		log.Errorln("error saving data to store: %s", err)
	}

	if _, ok := h.channels[channel]; !ok {
		// no one subscribed to this channel
		if m.Done != nil {
			m.Done <- nil
		}
		return
	}

	for conn := range h.channels[channel] {
		conn.Send(channel, data)
	}

	if m.Done != nil {
		m.Done <- nil
	}
}

type Stats struct {
	LocalSubscriptions map[string]int
}

func (h *Hub) Stats() (Stats, error) {
	h.Lock()
	defer h.Unlock()

	subscriptions := make(map[string]int)
	for k, v := range h.channels {
		subscriptions[k] = len(v)
	}

	return Stats{
		LocalSubscriptions: subscriptions,
	}, nil
}

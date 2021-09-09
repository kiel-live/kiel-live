package manager

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/manager/store"
	"github.com/kiel-live/kiel-live/protocol"
)

type Hub struct {
	subscriptions map[string]int
	client        *client.Client
	store         store.Store
	lock          sync.Mutex
}

func NewHub() *Hub {
	store := store.NewMemoryStore()
	store.Load()

	h := &Hub{
		subscriptions: make(map[string]int),
		store:         store,
	}

	return h
}

func (h *Hub) Unload() error {
	return h.store.Unload()
}

func (h *Hub) Subscribe(subject string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.subscriptions[subject]; !ok {
		// New channel!
		h.subscriptions[subject] = 1
		return nil
	}

	// TODO publish subscriptions list

	h.subscriptions[subject]++
	return nil
}

func (h *Hub) Unsubscribe(subject string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.subscriptions[subject]; !ok {
		return errors.New("No one subscribed to that subject")
	}

	h.subscriptions[subject] = h.subscriptions[subject] - 1

	// no one left for subject so delete it
	if h.subscriptions[subject] == 0 {
		delete(h.subscriptions, subject)
	}

	// TODO publish subscriptions list

	return nil
}

func (h *Hub) GetCache(subject string) (string, error) {
	return h.store.Get(subject)
}

func (h *Hub) SetCache(subject string, data string) error {
	return h.store.Set(subject, data)
}

func (h *Hub) updateSubscriptions() error {
	subscriptions := make([]string, len(h.subscriptions))

	i := 0
	for k := range h.subscriptions {
		subscriptions[i] = k
		i++
	}

	data, err := json.Marshal(subscriptions)
	if err != nil {
		return err
	}

	h.SetCache(protocol.SubjectSubscriptions, string(data))

	return h.client.Publish(protocol.SubjectSubscriptions, string(data))
}

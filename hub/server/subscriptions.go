package server

import (
	"encoding/json"
	"log/slog"
	"maps"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/pkg/metrics"
	"github.com/kiel-live/kiel-live/pkg/protocol"
)

// subscriptionManager tracks how many clients have subscribed to each detail topic
// and notifies registered collectors when the subscriber count crosses zero.
type subscriptionManager struct {
	mu         sync.Mutex
	counts     map[string]int
	collectors map[*collectorConn]struct{}
}

func newSubscriptionManager() *subscriptionManager {
	return &subscriptionManager{
		counts:     make(map[string]int),
		collectors: make(map[*collectorConn]struct{}),
	}
}

func (m *subscriptionManager) subscribe(topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	prev := m.counts[topic]
	m.counts[topic]++
	if prev == 0 {
		m.broadcastEvent(topic, true)
	}
}

func (m *subscriptionManager) unsubscribe(topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.counts[topic] <= 0 {
		return
	}
	m.counts[topic]--
	if m.counts[topic] == 0 {
		delete(m.counts, topic)
		m.broadcastEvent(topic, false)
	}
}

// broadcastEvent assumes mu is held.
func (m *subscriptionManager) broadcastEvent(topic string, subscribed bool) {
	env := buildSubscriptionEnvelope(topic, subscribed)
	for c := range m.collectors {
		c.sendRaw(env)
	}
}

func buildSubscriptionEnvelope(topic string, subscribed bool) []byte {
	se := protocol.SubscriptionEvent{Topic: topic, Subscribed: subscribed}
	data, _ := json.Marshal(se)
	raw := json.RawMessage(data)
	env := protocol.Envelope{
		Method: protocol.MethodUpdate,
		Topic:  protocol.TopicSubscriptions,
		Data:   &raw,
	}
	b, _ := json.Marshal(env)
	return b
}

func (m *subscriptionManager) registerCollector(c *collectorConn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectors[c] = struct{}{}
	// send current subscriptions as initial snapshot
	for topic := range m.counts {
		env := buildSubscriptionEnvelope(topic, true)
		c.sendRaw(env)
	}
}

func (m *subscriptionManager) unregisterCollector(c *collectorConn) {
	m.mu.Lock()
	delete(m.collectors, c)
	m.mu.Unlock()
}

func (m *subscriptionManager) registerMetrics(reg *metrics.Registry) {
	reg.Register("subscriptions.active", func() any {
		m.mu.Lock()
		defer m.mu.Unlock()
		return len(m.counts)
	})
	reg.Register("subscriptions.topics", func() any {
		m.mu.Lock()
		defer m.mu.Unlock()
		out := make(map[string]int, len(m.counts))
		maps.Copy(out, m.counts)
		return out
	})
}

// sendRaw sends a pre-marshalled frame to the collector; safe to call concurrently.
func (c *collectorConn) sendRaw(data []byte) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		slog.Debug("collector send error", "error", err)
	}
}

package pubsub

import "sync"

// PubSub is a thread-safe, in-memory publish-subscribe bus keyed by string topics.
// Multiple channels may subscribe to the same topic. Publish is non-blocking per
// subscriber: if a subscriber's channel is full the message is dropped for that
// subscriber.
type PubSub[T any] struct {
	mu   sync.RWMutex
	subs map[string]map[chan<- T]struct{}
}

func New[T any]() *PubSub[T] {
	return &PubSub[T]{subs: make(map[string]map[chan<- T]struct{})}
}

// Subscribe registers ch to receive messages published to topic.
// Calling Subscribe again with the same (topic, ch) pair is a no-op.
func (ps *PubSub[T]) Subscribe(topic string, ch chan<- T) {
	ps.mu.Lock()
	if ps.subs[topic] == nil {
		ps.subs[topic] = make(map[chan<- T]struct{})
	}
	ps.subs[topic][ch] = struct{}{}
	ps.mu.Unlock()
}

// Unsubscribe removes ch from topic. Safe to call even if not subscribed.
func (ps *PubSub[T]) Unsubscribe(topic string, ch chan<- T) {
	ps.mu.Lock()
	set := ps.subs[topic]
	delete(set, ch)
	if len(set) == 0 {
		delete(ps.subs, topic)
	}
	ps.mu.Unlock()
}

// Publish delivers msg to all current subscribers of topic.
// Per-subscriber delivery is non-blocking: slow subscribers are skipped.
func (ps *PubSub[T]) Publish(topic string, msg T) {
	ps.mu.RLock()
	for ch := range ps.subs[topic] {
		select {
		case ch <- msg:
		default:
		}
	}
	ps.mu.RUnlock()
}

// Subscribers returns the number of channels currently subscribed to topic.
func (ps *PubSub[T]) Subscribers(topic string) int {
	ps.mu.RLock()
	n := len(ps.subs[topic])
	ps.mu.RUnlock()
	return n
}

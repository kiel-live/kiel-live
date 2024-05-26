package pubsub

import (
	"context"
	"sync"
)

type Memory struct {
	sync.Mutex

	topics map[string]map[*Subscriber]struct{}
}

// New creates an in-memory publisher.
func NewMemory() Broker {
	return &Memory{
		topics: make(map[string]map[*Subscriber]struct{}),
	}
}

func (p *Memory) Publish(_ context.Context, _topic string, message Message) error {
	p.Lock()
	defer p.Unlock()

	topic, ok := p.topics[_topic]
	// if the topic does not exists we don't send anything
	if !ok {
		return nil
	}

	for s := range topic {
		go (*s)(message)
	}
	return nil
}

func (p *Memory) Subscribe(c context.Context, topic string, subscriber Subscriber) error {
	// Subscribe
	p.Lock()
	_, ok := p.topics[topic]
	if !ok {
		p.topics[topic] = make(map[*Subscriber]struct{})
	}
	p.topics[topic][&subscriber] = struct{}{}
	p.Unlock()

	// Wait for context to be done
	<-c.Done()

	// Unsubscribe
	p.Lock()
	delete(p.topics[topic], &subscriber)
	if len(p.topics[topic]) == 0 {
		delete(p.topics, topic)
	}
	p.Unlock()

	return nil
}

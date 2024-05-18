package pubsub_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kiel-live/kiel-live/hub/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestPubsub(t *testing.T) {
	var (
		wg sync.WaitGroup

		testTopic   = "world"
		testMessage = pubsub.Message("hello")
	)

	ctx, cancel := context.WithCancelCause(
		context.Background(),
	)

	broker := pubsub.NewMemory()
	go func() {
		broker.Subscribe(ctx, testTopic, func(message pubsub.Message) { assert.Equal(t, testMessage, message); wg.Done() })
	}()
	go func() {
		broker.Subscribe(ctx, testTopic, func(_ pubsub.Message) { wg.Done() })
	}()

	// Wait a bit for the subscriptions to be registered
	<-time.After(100 * time.Millisecond)

	wg.Add(2)
	go func() {
		broker.Publish(ctx, testTopic, testMessage)
	}()

	wg.Wait()
	cancel(nil)
}

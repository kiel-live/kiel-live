package pubsub_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kiel-live/kiel-live/pkg/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestPubsub(t *testing.T) {
	var (
		wg sync.WaitGroup

		testTopic   = "world"
		testMessage = pubsub.Message("hello")
	)

	ctx, cancel := context.WithCancelCause(t.Context())

	broker := pubsub.NewMemory()
	go func() {
		err := broker.Subscribe(ctx, testTopic, func(message pubsub.Message) {
			assert.Equal(t, testMessage, message)
			wg.Done()
		})
		assert.NoError(t, err)
	}()
	go func() {
		err := broker.Subscribe(ctx, testTopic, func(_ pubsub.Message) {
			wg.Done()
		})
		assert.NoError(t, err)
	}()

	// Wait a bit for the subscriptions to be registered
	<-time.After(100 * time.Millisecond)

	wg.Add(2)
	go func() {
		err := broker.Publish(ctx, testTopic, testMessage)
		assert.NoError(t, err)
	}()

	wg.Wait()
	cancel(nil)
}

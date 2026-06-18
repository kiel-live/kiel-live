package pubsub_test

import (
	"sync"
	"testing"
	"time"

	"github.com/kiel-live/kiel-live/pkg/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishDelivers(t *testing.T) {
	ps := pubsub.New[int]()
	ch := make(chan int, 1)
	ps.Subscribe("a", ch)

	ps.Publish("a", 42)

	select {
	case v := <-ch:
		assert.Equal(t, 42, v)
	case <-time.After(time.Second):
		t.Fatal("timeout: message not delivered")
	}
}

func TestPublishToUnsubscribedTopicIsNoop(t *testing.T) {
	ps := pubsub.New[string]()
	// no subscriber — must not panic
	ps.Publish("nothing", "hi")
}

func TestMultipleSubscribersSameTopic(t *testing.T) {
	ps := pubsub.New[string]()
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	ps.Subscribe("x", ch1)
	ps.Subscribe("x", ch2)

	ps.Publish("x", "hello")

	require.Equal(t, "hello", <-ch1)
	require.Equal(t, "hello", <-ch2)
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	ps := pubsub.New[int]()
	ch := make(chan int, 4)
	ps.Subscribe("t", ch)
	ps.Publish("t", 1)
	ps.Unsubscribe("t", ch)
	ps.Publish("t", 2) // must not reach ch

	require.Len(t, ch, 1, "only the pre-unsubscribe message should be in the channel")
	assert.Equal(t, 1, <-ch)
}

func TestUnsubscribeIdempotent(t *testing.T) {
	ps := pubsub.New[int]()
	ch := make(chan int, 1)
	ps.Subscribe("t", ch)
	ps.Unsubscribe("t", ch)
	ps.Unsubscribe("t", ch) // second call must not panic
}

func TestTopicsAreIndependent(t *testing.T) {
	ps := pubsub.New[string]()
	chA := make(chan string, 1)
	chB := make(chan string, 1)
	ps.Subscribe("a", chA)
	ps.Subscribe("b", chB)

	ps.Publish("a", "for-a")
	ps.Publish("b", "for-b")

	assert.Equal(t, "for-a", <-chA)
	assert.Equal(t, "for-b", <-chB)
	assert.Empty(t, chA, "b message must not reach a subscriber")
	assert.Empty(t, chB, "a message must not reach b subscriber")
}

func TestDropOnFullChannel(t *testing.T) {
	ps := pubsub.New[int]()
	ch := make(chan int, 1)
	ps.Subscribe("t", ch)

	ps.Publish("t", 1)
	ps.Publish("t", 2) // ch is full — must be dropped, not block

	assert.Len(t, ch, 1)
	assert.Equal(t, 1, <-ch)
}

func TestSubscribeAgainIsSafe(t *testing.T) {
	ps := pubsub.New[int]()
	ch := make(chan int, 2)
	ps.Subscribe("t", ch)
	ps.Subscribe("t", ch) // re-subscribe same channel

	ps.Publish("t", 7)

	// channel should have exactly one message (not duplicated)
	assert.Len(t, ch, 1)
	assert.Equal(t, 7, <-ch)
}

func TestSubscribers(t *testing.T) {
	ps := pubsub.New[int]()
	assert.Equal(t, 0, ps.Subscribers("t"))

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ps.Subscribe("t", ch1)
	assert.Equal(t, 1, ps.Subscribers("t"))
	ps.Subscribe("t", ch2)
	assert.Equal(t, 2, ps.Subscribers("t"))

	ps.Unsubscribe("t", ch1)
	assert.Equal(t, 1, ps.Subscribers("t"))
	ps.Unsubscribe("t", ch2)
	assert.Equal(t, 0, ps.Subscribers("t"))
}

func TestConcurrentPublishSubscribe(t *testing.T) {
	ps := pubsub.New[int]()
	const workers = 8
	const msgs = 100

	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ch := make(chan int, msgs)
			topic := "topic"
			ps.Subscribe(topic, ch)
			for j := range msgs {
				ps.Publish(topic, id*1000+j)
			}
			ps.Unsubscribe(topic, ch)
		}(i)
	}
	wg.Wait() // must complete without race or deadlock
}

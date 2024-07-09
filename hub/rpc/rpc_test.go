package rpc_test

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiel-live/kiel-live/hub/rpc"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

type SampleRPC struct {
}

func (t *SampleRPC) Hello(name string) (string, error) {
	return "Hello, " + name, nil
}

func BenchmarkRPC(b *testing.B) {
	ctx := context.Background()

	serverPeer, clientPeer := net.Pipe()

	broker := pubsub.NewMemory()

	server := rpc.NewServerPeer(ctx, serverPeer, broker)
	assert.NotNil(b, server)
	err := server.Register(&SampleRPC{})
	assert.NoError(b, err)

	client := rpc.NewClientPeer(ctx, clientPeer)

	args := []any{"Alice"}
	response := []any{}
	err = client.Request(ctx, "Hello", args, &response)
	assert.NoError(b, err)
	assert.Equal(b, []any{"Hello, Alice"}, response)
}

func BenchmarkSubscribe(b *testing.B) {
	ctx := context.Background()

	serverPeer, clientPeer := net.Pipe()

	broker := pubsub.NewMemory()

	server := rpc.NewServerPeer(ctx, serverPeer, broker)
	assert.NotNil(b, server)

	client := rpc.NewClientPeer(ctx, clientPeer)

	g := sync.WaitGroup{}

	channelName := "test-channel"

	g.Add(1)
	err := client.Subscribe(ctx, channelName, func(a any) {
		assert.Equal(b, "Hello, World", a)
		g.Done()
	})
	assert.NoError(b, err)

	g.Add(1)
	go func() {
		err = broker.Publish(context.Background(), channelName, []byte("Hello, World"))
		assert.NoError(b, err)

		// err = client.Publish(ctx, channelName, []byte("Hello, World from the client"))
		// assert.NoError(b, err)
		g.Done()
	}()

	g.Wait()

	err = client.Unsubscribe(ctx, channelName)
	assert.NoError(b, err)
}

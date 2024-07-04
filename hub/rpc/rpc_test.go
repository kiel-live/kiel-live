package rpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiel-live/kiel-live/hub/rpc"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

type SampleRPC struct {
}

type HelloArgs struct {
	Name string `json:"name"`
}

func (t *SampleRPC) Hello(args *HelloArgs, reply *string) error {
	*reply = "Hello, " + args.Name
	return nil
}

func BenchmarkPeer(b *testing.B) {
	serverPeer, clientPeer := net.Pipe()

	broker := pubsub.NewMemory()

	_, err := rpc.NewServer(&SampleRPC{}, broker, serverPeer)
	assert.NoError(b, err)

	client := rpc.NewClient(clientPeer)

	args := &HelloArgs{Name: "Alice"}
	var response string
	err = client.Request("Hello", args, &response)
	assert.NoError(b, err)
	assert.Equal(b, "Hello, Alice", response)

	err = client.Subscribe("test-channel")
	assert.NoError(b, err)

	err = broker.Publish(context.Background(), "test-channel", []byte("Hello, World"))
	assert.NoError(b, err)

	err = client.Publish("test-channel", []byte("Hello, World from the client"))
	assert.NoError(b, err)

	err = client.Unsubscribe("test-channel")
	assert.NoError(b, err)
}

package rpc_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/stretchr/testify/assert"

	"github.com/kiel-live/kiel-live/pkg/pubsub"
	"github.com/kiel-live/kiel-live/pkg/rpc"
)

type SampleRPC struct{}

func (t *SampleRPC) Hello(name string) (string, error) {
	return "Hello, " + name, nil
}

func ProxyConnection(clientConn net.Conn, serverConn net.Conn) {
	defer clientConn.Close()
	defer serverConn.Close()

	// Channel to signal completion
	done := make(chan struct{})

	// Log and forward messages from client to server
	go func() {
		clientScanner := bufio.NewScanner(clientConn)
		for clientScanner.Scan() {
			message := clientScanner.Text()
			fmt.Printf("Client to Server: %s\n", message)
			_, err := fmt.Fprintln(serverConn, message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error forwarding to server: %v\n", err)
				break
			}
		}
		done <- struct{}{}
	}()

	// Log and forward messages from server to client
	go func() {
		serverScanner := bufio.NewScanner(serverConn)
		for serverScanner.Scan() {
			message := serverScanner.Text()
			fmt.Printf("Server to Client: %s\n", message)
			_, err := fmt.Fprintln(clientConn, message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error forwarding to client: %v\n", err)
				break
			}
		}
		done <- struct{}{}
	}()

	<-done
}

func BenchmarkRPC(b *testing.B) {
	ctx := context.Background()

	// serverPeer, proxyPeerServer := net.Pipe()
	// proxyPeerClient, clientPeer := net.Pipe()
	// go ProxyConnection(proxyPeerClient, proxyPeerServer)

	serverPeer, clientPeer := net.Pipe()

	broker := pubsub.NewMemory()

	server := rpc.NewServer(broker)
	assert.NotNil(b, server)
	err := server.Register(&SampleRPC{})

	sP := server.NewPeer(ctx, jsonrpc2.NewPlainObjectStream(serverPeer))
	defer sP.Close()
	assert.NotNil(b, sP)

	client := rpc.NewClientPeer(ctx, jsonrpc2.NewPlainObjectStream(clientPeer))

	args := []any{"Alice"}
	response := []any{}
	err = client.Call(ctx, "Hello", args, &response)
	assert.NoError(b, err)
	assert.Equal(b, []any{"Hello, Alice"}, response)
}

func TestSubscribe(t *testing.T) {
	ctx := context.Background()

	serverPeer, proxyPeerServer := net.Pipe()
	proxyPeerClient, clientPeer := net.Pipe()
	go ProxyConnection(proxyPeerClient, proxyPeerServer)

	broker := pubsub.NewMemory()
	server := rpc.NewServer(broker)
	assert.NotNil(t, server)

	sP := server.NewPeer(ctx, jsonrpc2.NewPlainObjectStream(serverPeer))
	defer sP.Close()
	assert.NotNil(t, sP)

	client := rpc.NewClientPeer(ctx, jsonrpc2.NewPlainObjectStream(clientPeer))

	g := sync.WaitGroup{}

	channelName := "test-channel"

	g.Add(1)
	err := client.Subscribe(ctx, channelName, func(a *json.RawMessage) {
		assert.JSONEq(t, `"Hello, World from the client"`, string(*a))
		g.Done()
	})
	assert.NoError(t, err)

	g.Add(1)
	go func() {
		// err = broker.Publish(context.Background(), channelName, []byte("Hello, World"))
		// assert.NoError(b, err)

		err = client.Publish(ctx, channelName, "Hello, World from the client")
		assert.NoError(t, err)
		g.Done()
	}()

	g.Wait()

	err = client.Unsubscribe(ctx, channelName)
	assert.NoError(t, err)
}

func TestLol(t *testing.T) {
	_msg := json.RawMessage("test")
	msg := &_msg

	d := []byte(*msg)
	assert.Equal(t, []byte("test"), d)
}

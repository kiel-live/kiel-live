package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

type ClientPeer struct {
	defaultServiceName string
	client             *jsonrpc2.Conn
	mutex              sync.Mutex
	subscriptions      map[string]Subscription
}

func NewClientPeer(ctx context.Context, conn io.ReadWriteCloser) *ClientPeer {
	peer := &ClientPeer{
		defaultServiceName: "main",
		subscriptions:      make(map[string]Subscription),
	}

	rpcConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(conn), peer)
	peer.client = rpcConn

	return peer
}

func (c *ClientPeer) handlePubsubMessage(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	c.mutex.Lock()
	sub, exists := c.subscriptions[r.Method]
	c.mutex.Unlock()

	if exists {
		var data string // TODO: reflect correct data type
		if err := json.Unmarshal(*r.Params, &data); err != nil {
			return err
		}

		sub(data)
		return nil
	}

	return nil
}

// Handle implements the jsonrpc2.Handler interface.
func (c *ClientPeer) Handle(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) {
	// check if this is a pubsub notification
	if r.Notif {
		c.handlePubsubMessage(ctx, conn, r)
	}

	// we don't care about requests as we are only a client
}

func (c *ClientPeer) Request(ctx context.Context, method string, args any, reply any) error {
	return c.client.Call(ctx, fmt.Sprintf("%s.%s", c.defaultServiceName, method), args, reply)
}

func (c *ClientPeer) Subscribe(ctx context.Context, channel string, sub func(any)) error {
	// check if we are already subscribed to this channel
	c.mutex.Lock()
	_, exists := c.subscriptions[channel]
	c.mutex.Unlock()
	if exists {
		return fmt.Errorf("already subscribed to channel: %s", channel)
	}

	request := &SubscribeRequest{
		Channel: channel,
	}

	var response string
	err := c.client.Call(ctx, fmt.Sprintf("%s.Subscribe", internalServiceName), request, &response)
	if err != nil {
		return err
	}

	if response != "ok" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	c.mutex.Lock()
	c.subscriptions[channel] = sub
	c.mutex.Unlock()

	return nil
}

func (c *ClientPeer) Unsubscribe(ctx context.Context, channel string) error {
	// check if we are subscribed to this channel
	c.mutex.Lock()
	_, exists := c.subscriptions[channel]
	c.mutex.Unlock()
	if !exists {
		return fmt.Errorf("not subscribed to channel: %s", channel)
	}

	request := &UnsubscribeRequest{
		Channel: channel,
	}

	var response string
	err := c.client.Call(ctx, fmt.Sprintf("%s.Unsubscribe", internalServiceName), request, &response)
	if err != nil {
		return err
	}

	if response != "ok" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	c.mutex.Lock()
	delete(c.subscriptions, channel)
	c.mutex.Unlock()
	// if p.subscriptions[channel] == 0 {

	return nil
}

func (c *ClientPeer) Publish(ctx context.Context, channel string, _data any) error {
	data, err := json.Marshal(_data)
	if err != nil {
		return err
	}

	request := &PublishRequest{
		Channel: channel,
		Data:    data,
	}

	// TODO: should we care about a response?
	// var response string
	err = c.client.Notify(ctx, fmt.Sprintf("%s.Publish", internalServiceName), request)
	if err != nil {
		return err
	}

	return nil
}

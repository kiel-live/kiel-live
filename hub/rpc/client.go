package rpc

import (
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Client struct {
	client        *rpc.Client
	subscriptions map[string]int
}

func NewClient(conn io.ReadWriteCloser) *Client {
	peer := &Client{
		client:        jsonrpc.NewClient(conn),
		subscriptions: make(map[string]int),
	}

	return peer
}

func (p *Client) Request(method string, args any, reply any) error {
	return p.client.Call(fmt.Sprintf("%s.%s", serviceName, method), args, reply)
}

func (p *Client) Subscribe(channel string) error {
	// check if we are already subscribed to this channel
	if _, exists := p.subscriptions[channel]; exists {
		p.subscriptions[channel]++
		return nil
	}

	request := &SubscribeRequest{
		Channel: channel,
	}

	var response string
	err := p.client.Call(fmt.Sprintf("%s.Subscribe", internalServiceName), request, &response)
	if err != nil {
		return err
	}

	if response != "ok" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	p.subscriptions[channel]++

	return nil
}

func (p *Client) Unsubscribe(channel string) error {
	// check if we are subscribed to this channel
	if _, exists := p.subscriptions[channel]; !exists {
		return fmt.Errorf("not subscribed to channel: %s", channel)
	}

	request := &UnsubscribeRequest{
		Channel: channel,
	}

	var response string
	err := p.client.Call(fmt.Sprintf("%s.Unsubscribe", internalServiceName), request, &response)
	if err != nil {
		return err
	}

	if response != "ok" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	p.subscriptions[channel]--
	if p.subscriptions[channel] == 0 {
		delete(p.subscriptions, channel)
	}

	return nil
}

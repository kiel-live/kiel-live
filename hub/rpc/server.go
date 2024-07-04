package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"

	"github.com/kiel-live/kiel-live/shared/pubsub"
)

type Server struct {
	*rpc.Server
}

func (s *Server) Publish(channel string, _data any) error {
	data, err := json.Marshal(_data)
	if err != nil {
		return err
	}

	_msg := &ChannelMessage{
		Channel: channel,
		Data:    data,
	}
	msg, err := json.Marshal(_msg)
	if err != nil {
		return err
	}

	// TODO: broadcast to all subscribers
	fmt.Println("Publish", string(msg))
	// return s.Broadcast("Publish", string(msg))

	return nil
}

type SubscriptionsRPC struct {
	sync.Mutex
	client        io.ReadWriteCloser
	pubsub        pubsub.Broker
	subscriptions map[string]context.Context
}

func (s *SubscriptionsRPC) Subscribe(args *SubscribeRequest, reply *string) error {
	fmt.Println("Subscribe", args.Channel)

	s.Lock()
	if _, exists := s.subscriptions[args.Channel]; exists {
		*reply = "already subscribed"
		return nil
	}
	s.Unlock()

	ctx := context.Background()
	err := s.pubsub.Subscribe(ctx, args.Channel, func(message pubsub.Message) {
		_msg := &ChannelMessage{
			Channel: args.Channel,
			Data:    message,
		}
		msg, err := json.Marshal(_msg)
		if err != nil {
			return
		}

		fmt.Println("Forward message", string(msg), "to client")

		_, _ = s.client.Write(msg)
	})
	if err != nil {
		*reply = err.Error()
		return err
	}

	s.Lock()
	s.subscriptions[args.Channel] = ctx
	s.Unlock()

	*reply = "ok"
	return nil
}

func (s *SubscriptionsRPC) Unsubscribe(args *UnsubscribeRequest, reply *string) error {
	s.Lock()
	defer s.Unlock()

	fmt.Println("Unsubscribe", args.Channel)

	if _, exists := s.subscriptions[args.Channel]; !exists {
		*reply = "not subscribed"
		return fmt.Errorf("not subscribed to channel: %s", args.Channel)
	}

	s.subscriptions[args.Channel].Done()
	*reply = "ok"
	return nil
}

func NewServer(rcvr any, pubsub pubsub.Broker, client io.ReadWriteCloser) (*Server, error) {
	server := &Server{
		rpc.NewServer(),
	}

	err := server.RegisterName(serviceName, rcvr)
	if err != nil {
		return nil, err
	}

	subscriptions := &SubscriptionsRPC{
		pubsub:        pubsub, // TODO
		client:        client,
		subscriptions: make(map[string]context.Context),
	}

	err = server.RegisterName(internalServiceName, subscriptions)
	if err != nil {
		return nil, err
	}

	go server.ServeCodec(jsonrpc.NewServerCodec(client))

	return server, nil
}

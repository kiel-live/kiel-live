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
		return fmt.Errorf("already subscribed to channel: %s", args.Channel)
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
		return fmt.Errorf("not subscribed to channel: %s", args.Channel)
	}

	s.subscriptions[args.Channel].Done()
	*reply = "ok"
	return nil
}

func (s *SubscriptionsRPC) Publish(args *PublishRequest, reply *string) error {
	fmt.Println("Publish", args.Channel)

	ctx := context.Background()
	err := s.pubsub.Publish(ctx, args.Channel, args.Data)
	if err != nil {
		return err
	}

	s.subscriptions[args.Channel].Done()
	*reply = "ok"
	return nil
}

func NewServer(rcvr any, pubsub pubsub.Broker, client io.ReadWriteCloser) error {
	server := rpc.NewServer()

	err := server.RegisterName(serviceName, rcvr)
	if err != nil {
		return err
	}

	subscriptions := &SubscriptionsRPC{
		pubsub:        pubsub, // TODO
		client:        client,
		subscriptions: make(map[string]context.Context),
	}

	err = server.RegisterName(internalServiceName, subscriptions)
	if err != nil {
		return err
	}

	go server.ServeCodec(jsonrpc.NewServerCodec(client))

	return nil
}

package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
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
}

func (t *SubscriptionsRPC) Subscribe(args *SubscribeRequest, reply *string) error {
	fmt.Println("Subscribe", args.Channel)
	*reply = "ok"
	return nil
}

func (t *SubscriptionsRPC) Unsubscribe(args *UnsubscribeRequest, reply *string) error {
	fmt.Println("Unsubscribe", args.Channel)
	*reply = "ok"
	return nil
}

func NewServer(rcvr any, conn io.ReadWriteCloser) (*Server, error) {
	server := &Server{
		rpc.NewServer(),
	}

	err := server.RegisterName(serviceName, rcvr)
	if err != nil {
		return nil, err
	}

	subscriptions := &SubscriptionsRPC{}

	err = server.RegisterName(internalServiceName, subscriptions)
	if err != nil {
		return nil, err
	}

	go server.ServeCodec(jsonrpc.NewServerCodec(conn))

	return server, nil
}

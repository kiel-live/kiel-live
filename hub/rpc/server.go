package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/kiel-live/kiel-live/hub/rpc/service"
	"github.com/kiel-live/kiel-live/shared/pubsub"
	"github.com/sourcegraph/jsonrpc2"
)

type Handler interface {
	Handle(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request)
}

type ServerPeer struct {
	services      map[string]*service.Service
	client        *jsonrpc2.Conn
	mutex         sync.Mutex
	broker        pubsub.Broker
	subscriptions map[*jsonrpc2.Conn]map[string]context.Context
}

func NewServerPeer(ctx context.Context, conn io.ReadWriteCloser, broker pubsub.Broker) *ServerPeer {
	peer := &ServerPeer{
		services:      make(map[string]*service.Service),
		broker:        broker,
		subscriptions: make(map[*jsonrpc2.Conn]map[string]context.Context),
	}

	rpcConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(conn), peer)
	peer.client = rpcConn

	return peer
}

func (p *ServerPeer) RegisterName(name string, st any) error {
	if _, exists := p.services[name]; exists {
		return fmt.Errorf("service already registered: %s", name)
	}

	s, err := service.NewService(st)
	if err != nil {
		return err
	}

	p.services[name] = s
	return nil
}

func (p *ServerPeer) Register(s any) error {
	return p.RegisterName(defaultServiceName, s)
}

// Handle implements the jsonrpc2.Handler interface.
func (p *ServerPeer) Handle(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) {
	if r.Notif {
		err := p.handlePublish(ctx, conn, r)
		if err != nil {
			err := p.client.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
				Code:    jsonrpc2.CodeInternalError,
				Message: err.Error(),
			})
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	err := p.handleWithError(ctx, conn, r)
	if err != nil {
		err := p.client.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: err.Error(),
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *ServerPeer) handleWithError(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	switch r.Method {
	case fmt.Sprintf("%s.Subscribe", internalServiceName):
		return p.handleSubscribe(ctx, conn, r)

	case fmt.Sprintf("%s.Unsubscribe", internalServiceName):
		return p.handleUnsubscribe(ctx, conn, r)
	}

	// find service
	serviceName, methodName, found := strings.Cut(r.Method, ".")
	if !found {
		return fmt.Errorf("invalid method name")
	}

	s, exists := p.services[serviceName]
	if !exists {
		return fmt.Errorf("service not found")
	}

	// find method
	results, err := s.Call(methodName, r.Params)
	if err != nil {
		return err
	}

	return p.client.Reply(ctx, r.ID, results)
}

func (p *ServerPeer) handlePublish(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	// TODO: check auth

	// forward message to broker
	d := []byte(*r.Params)
	return p.broker.Publish(ctx, r.Method, d)
}

func (p *ServerPeer) handleSubscribe(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	var args SubscribeRequest
	if err := json.Unmarshal(*r.Params, &args); err != nil {
		return err
	}

	p.mutex.Lock()
	subs, exists := p.subscriptions[conn]
	if !exists {
		subs = make(map[string]context.Context)
		p.subscriptions[conn] = subs
	}
	p.mutex.Unlock()

	if _, exists := subs[args.Channel]; exists {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("already subscribed to channel: %s", args.Channel),
		})
	}

	subCtx := context.Background()

	err := p.broker.Subscribe(subCtx, args.Channel, func(_message pubsub.Message) {
		message := json.RawMessage(_message)
		if err := conn.Notify(ctx, args.Channel, message); err != nil {
			log.Println(err)
			return
		}
	})

	if err != nil {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: err.Error(),
		})
	}

	p.mutex.Lock()
	subs[args.Channel] = subCtx
	p.subscriptions[conn] = subs
	p.mutex.Unlock()

	return conn.Reply(ctx, r.ID, "ok")
}

func (p *ServerPeer) handleUnsubscribe(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	var args UnsubscribeRequest
	if err := json.Unmarshal(*r.Params, &args); err != nil {
		return err
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	subs, exists := p.subscriptions[conn]

	if !exists {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("not subscribed to channel: %s", args.Channel),
		})
	}

	subCtx, exists := subs[args.Channel]
	if !exists {
		return conn.ReplyWithError(subCtx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("not subscribed to channel: %s", args.Channel),
		})
	}

	subCtx.Done()
	delete(subs, args.Channel)

	if len(subs) == 0 {
		delete(p.subscriptions, conn)
	}

	return conn.Reply(ctx, r.ID, "ok")
}

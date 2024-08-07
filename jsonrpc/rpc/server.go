package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sourcegraph/jsonrpc2"

	"github.com/kiel-live/kiel-live/jsonrpc/rpc/service"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

type Handler interface {
	Handle(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request)
}

type Server struct {
	services      map[string]*service.Service
	broker        pubsub.Broker
	mutex         sync.Mutex
	subscriptions map[*jsonrpc2.Conn]map[string]context.CancelFunc
}

func NewServer(broker pubsub.Broker) *Server {
	return &Server{
		services:      make(map[string]*service.Service),
		broker:        broker,
		mutex:         sync.Mutex{},
		subscriptions: make(map[*jsonrpc2.Conn]map[string]context.CancelFunc),
	}
}

func (s *Server) RegisterName(name string, st any) error {
	if _, exists := s.services[name]; exists {
		return fmt.Errorf("service already registered: %s", name)
	}

	srv, err := service.NewService(st)
	if err != nil {
		return err
	}

	s.services[name] = srv
	return nil
}

func (s *Server) Register(srv any) error {
	return s.RegisterName(defaultServiceName, srv)
}

func (s *Server) getSubscriptions(conn *jsonrpc2.Conn) map[string]context.CancelFunc {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	subs, exists := s.subscriptions[conn]
	if !exists {
		return nil
	}

	return subs
}

func (s *Server) addSubscription(conn *jsonrpc2.Conn, channel string, unsubscribe context.CancelFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subs, exists := s.subscriptions[conn]
	if !exists {
		subs = make(map[string]context.CancelFunc)
		s.subscriptions[conn] = subs
	}

	subs[channel] = unsubscribe
}

func (s *Server) removeSubscription(conn *jsonrpc2.Conn, channel string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subs, exists := s.subscriptions[conn]
	if !exists {
		return false
	}

	delete(subs, channel)

	if len(subs) == 0 {
		delete(s.subscriptions, conn)
	}

	return true
}

func (s *Server) NewPeer(ctx context.Context, conn jsonrpc2.ObjectStream) *jsonrpc2.Conn {
	return jsonrpc2.NewConn(ctx, conn, s)
}

// Handle implements the jsonrpc2.Handler interface.
func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) {
	// fmt.Printf("%s: %d\n", "handle", time.Now().UnixMicro())

	if r.Notif {
		err := s.handlePublish(ctx, conn, r)
		if err != nil {
			err := conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
				Code:    jsonrpc2.CodeInternalError,
				Message: err.Error(),
			})
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	err := s.handleWithError(ctx, conn, r)
	if err != nil {
		err := conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: err.Error(),
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) handleWithError(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	switch r.Method {
	case fmt.Sprintf("%s.Subscribe", internalServiceName):
		return s.handleSubscribe(ctx, conn, r)

	case fmt.Sprintf("%s.Unsubscribe", internalServiceName):
		return s.handleUnsubscribe(ctx, conn, r)
	}

	// find service
	serviceName, methodName, found := strings.Cut(r.Method, ".")
	if !found {
		return fmt.Errorf("invalid method name")
	}

	srv, exists := s.services[serviceName]
	if !exists {
		return fmt.Errorf("service not found")
	}

	// find method
	results, err := srv.Call(methodName, r.Params)
	if err != nil {
		return err
	}

	return conn.Reply(ctx, r.ID, results)
}

func (s *Server) handlePublish(ctx context.Context, _ *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	// TODO: check auth

	// forward message to broker
	// fmt.Printf("%s: %d\n", "handle publish", time.Now().UnixMicro())
	d := []byte(*r.Params)
	return s.broker.Publish(ctx, r.Method, d)
}

func (s *Server) handleSubscribe(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	var args SubscribeRequest
	if err := json.Unmarshal(*r.Params, &args); err != nil {
		return err
	}

	subs := s.getSubscriptions(conn)
	if subs == nil {
		subs = make(map[string]context.CancelFunc)
	}

	if _, exists := subs[args.Channel]; exists {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("already subscribed to channel: %s", args.Channel),
		})
	}

	subCtx, unsubscribe := context.WithCancel(context.Background())

	err := s.broker.Subscribe(subCtx, args.Channel, func(_message pubsub.Message) {
		// fmt.Printf("%s: %d\n", "forward msg", time.Now().UnixMicro())
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

	s.addSubscription(conn, args.Channel, unsubscribe)

	// fmt.Printf("%s: %d\n", "subscribed", time.Now().UnixMicro())

	return conn.Reply(ctx, r.ID, "ok")
}

func (s *Server) handleUnsubscribe(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) error {
	var args UnsubscribeRequest
	if err := json.Unmarshal(*r.Params, &args); err != nil {
		return err
	}

	subs := s.getSubscriptions(conn)
	if subs == nil {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("not subscribed to channel: %s", args.Channel),
		})
	}

	unsubscribe, exists := subs[args.Channel]
	if !exists {
		return conn.ReplyWithError(ctx, r.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("not subscribed to channel: %s", args.Channel),
		})
	}

	unsubscribe()

	s.removeSubscription(conn, args.Channel)

	// fmt.Printf("%s: %d\n", "unsubscribed", time.Now().UnixMicro())

	return conn.Reply(ctx, r.ID, "ok")
}

func (s *Server) ClosePeer(conn *jsonrpc2.Conn) error {
	subs := s.getSubscriptions(conn)
	for channel, unsubscribe := range subs {
		unsubscribe()
		s.removeSubscription(conn, channel)
	}
	return conn.Close()
}

func (s *Server) Close() error {
	for conn := range s.subscriptions {
		err := s.ClosePeer(conn)
		if err != nil {
			return err
		}
	}

	return nil
}

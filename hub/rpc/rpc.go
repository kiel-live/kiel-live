package rpc

import "github.com/kiel-live/kiel-live/shared/hub"

type RPC struct {
	hub *hub.Hub
}

func NewRPC(hub *hub.Hub) *RPC {
	return &RPC{
		hub: hub,
	}
}

type HelloArgs struct {
	Name string `json:"name"`
}

func (t *RPC) Hello(args *HelloArgs, reply *string) error {
	*reply = "Hello, " + args.Name
	return nil
}

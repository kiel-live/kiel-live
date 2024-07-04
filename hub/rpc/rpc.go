package rpc

import "github.com/kiel-live/kiel-live/shared/hub"

type RPC struct {
	hub *hub.Hub
}

type HelloArgs struct {
	Name string `json:"name"`
}

func (t *RPC) Hello(args *HelloArgs, reply *string) error {
	*reply = "Hello, " + args.Name
	return nil
}

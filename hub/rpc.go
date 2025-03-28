package main

import (
	"github.com/kiel-live/kiel-live/pkg/hub"
)

type KielLiveRPC struct {
	Hub *hub.Hub
}

func (t *KielLiveRPC) Hello(name string) (string, error) {
	return "Hello, " + name, nil
}

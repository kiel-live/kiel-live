package graph

import (
	"fmt"
	"strings"

	"github.com/kiel-live/kiel-live/hub/graph/model"
	"github.com/tidwall/buntdb"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *buntdb.DB
	Channels map[string]chan *model.Map
}

func (r *Resolver) OpenMapChannel(id string) chan *model.Map {
	ch := make(chan *model.Map, 1)
	r.Channels[id] = ch
	return ch
}

func (r *Resolver) GetMapChannel(id string) (chan *model.Map, bool) {
	ch, ok := r.Channels[id]
	return ch, ok
}

func (r *Resolver) CloseMapChannel(id string) {
	close(r.Channels[id])
	delete(r.Channels, id)
}

func (r *Resolver) GetMapChannels(lat, lng float64) (map[string]chan *model.Map, error) {
	channels := make(map[string]chan *model.Map)

	err := r.DB.View(func(tx *buntdb.Tx) error {
		pos := fmt.Sprintf("[%f %f]", lat, lng)
		return tx.Intersects("subscription_map", pos, func(key, val string) bool {
			channelID := strings.Replace(key, "subscription:map:", "", 1)
			ch, ok := r.GetMapChannel(channelID)
			if ok {
				channels[channelID] = ch
			}
			return true
		})
	})

	return channels, err
}

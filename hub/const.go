package main

import (
	"context"
	"fmt"
)

const (
	MapItemUpdatedTopic = "map.updated:%s-%s"
	MapItemDeletedTopic = "map.deleted:%s-%s"
	ItemUpdatedTopic    = "item.updated:%s-%s"
	ItemDeletedTopic    = "item.deleted:%s-%s"
)

type Model interface {
	ID() string
	ToJSON() []byte
}

type MapModel interface {
	*Model
	Location() string
}

func (h *Hub) updateMapItem(ctx context.Context, name string, model Model) error {
	return h.PubSub.Publish(ctx, fmt.Sprintf(MapItemUpdatedTopic, name, model.ID()), model.ToJSON())
}

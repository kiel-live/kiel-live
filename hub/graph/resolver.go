package graph

import (
	"context"
	"fmt"

	"github.com/kiel-live/kiel-live/shared/hub"
	"github.com/kiel-live/kiel-live/shared/models"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Hub *hub.Hub
}

func (r *Resolver) subscribeBoundingBox(ctx context.Context, minLat float64, minLng float64, maxLat float64, maxLng float64, topicPrefix string, subscriber func(pubsub.Message)) error {
	cellIDs := (&models.BoundingBox{
		MinLat: minLat,
		MinLng: minLng,
		MaxLat: maxLat,
		MaxLng: maxLng,
	}).GetCellIDs()

	for _, cellID := range cellIDs {
		err := r.Hub.PubSub.Subscribe(ctx, fmt.Sprintf("%s:%d", topicPrefix, cellID), subscriber)
		if err != nil {
			return err
		}
	}

	return nil
}

package hub

import (
	"context"
	"fmt"

	"github.com/kiel-live/kiel-live/hub/graph/model"
	"github.com/kiel-live/kiel-live/hub/pubsub"
)

type Hub struct {
	pubsub pubsub.Broker
}

func (h *Hub) SetVehicle(vehicle *model.Vehicle) {
	ctx := context.Background()
	h.pubsub.Publish(ctx, fmt.Sprintf("vehicle:%s", vehicle.ID), []byte{})

	// TODO: publish to all cell levels
	h.pubsub.Publish(ctx, fmt.Sprintf("map:%s", vehicle.Location.GetCellID().ToToken()), []byte{})
}

func (h *Hub) SetStop(stop *model.Stop) {
	ctx := context.Background()
	h.pubsub.Publish(ctx, fmt.Sprintf("vehicle:%s", stop.ID), []byte{})

	// TODO: publish to all cell levels
	h.pubsub.Publish(ctx, fmt.Sprintf("map:%s", stop.Location.GetCellID().ToToken()), []byte{})
}

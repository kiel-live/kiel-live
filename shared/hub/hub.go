package hub

import (
	"context"
	"fmt"

	"github.com/kiel-live/kiel-live/shared/models"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

type Hub struct {
	pubsub pubsub.Broker
}

func (h *Hub) SetVehicle(vehicle *models.Vehicle) error {
	ctx := context.Background()
	err := h.pubsub.Publish(ctx, fmt.Sprintf("vehicle:%s", vehicle.ID), []byte{})
	if err != nil {
		return err
	}

	cellIDs := vehicle.Location.GetCellIDs()
	for _, cellID := range cellIDs {
		err = h.pubsub.Publish(ctx, fmt.Sprintf("map:%s", cellID.ToToken()), []byte{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) SetStop(stop *models.Stop) error {
	ctx := context.Background()
	err := h.pubsub.Publish(ctx, fmt.Sprintf("vehicle:%s", stop.ID), []byte{})
	if err != nil {
		return err
	}

	cellIDs := stop.Location.GetCellIDs()
	for _, cellID := range cellIDs {
		err = h.pubsub.Publish(ctx, fmt.Sprintf("map:%s", cellID.ToToken()), []byte{})
		if err != nil {
			return err
		}
	}

	return nil
}

package main

import (
	"context"
	"fmt"

	"github.com/kiel-live/kiel-live/pkg/database"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/pubsub"
)

// TODO: use functions as interface for client as well
type Hub struct {
	DB     database.Database
	PubSub pubsub.Broker
}

func (h *Hub) GetVehicle(ctx context.Context, vehicleID string) (*models.Vehicle, error) {
	return h.DB.GetVehicle(ctx, vehicleID)
}

func (h *Hub) GetVehicles(ctx context.Context, opts *database.ListOptions) ([]*models.Vehicle, error) {
	return h.DB.GetVehicles(ctx, opts)
}

func (h *Hub) SetVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	err := h.DB.SetVehicle(ctx, vehicle)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("vehicle-updated:%s", vehicle.ID), vehicle.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range vehicle.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf("map-vehicle-updated:%s", cellID.ToToken()), vehicle.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) DeleteVehicle(ctx context.Context, vehicleID string) error {
	vehicle, err := h.DB.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	err = h.DB.DeleteVehicle(ctx, vehicle.ID)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("vehicle-deleted:%s", vehicle.ID), vehicle.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range vehicle.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf("map-vehicle-deleted:%s", cellID.ToToken()), vehicle.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) GetStop(ctx context.Context, stopID string) (*models.Stop, error) {
	return h.DB.GetStop(ctx, stopID)
}

func (h *Hub) GetStops(ctx context.Context, opts *database.ListOptions) ([]*models.Stop, error) {
	return h.DB.GetStops(ctx, opts)
}

func (h *Hub) SetStop(ctx context.Context, stop *models.Stop) error {
	err := h.DB.SetStop(ctx, stop)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("stop-updated:%s", stop.ID), stop.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range stop.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf("map-stop-updated:%s", cellID.ToToken()), stop.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) DeleteStop(ctx context.Context, stopID string) error {
	stop, err := h.DB.GetStop(ctx, stopID)
	if err != nil {
		return err
	}

	err = h.DB.DeleteStop(ctx, stop.ID)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("stop-deleted:%s", stop.ID), stop.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range stop.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf("map-stop-deleted:%s", cellID.ToToken()), stop.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

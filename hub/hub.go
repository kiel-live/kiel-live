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
	// TODO: test auth

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
	// TODO: test auth

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
	// TODO: test auth

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
	// TODO: test auth

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

func (h *Hub) GetTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	return h.DB.GetTrip(ctx, tripID)
}

func (h *Hub) SetTrip(ctx context.Context, trip *models.Trip) error {
	// TODO: require auth
	err := h.DB.SetTrip(ctx, trip)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("trip-updated:%s", trip.ID), trip.ToJSON())
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) GetRoute(ctx context.Context, tripID string) (*models.Route, error) {
	return h.DB.GetRoute(ctx, tripID)
}

func (h *Hub) SetRoute(ctx context.Context, route *models.Route) error {
	// TODO: require auth
	err := h.DB.SetRoute(ctx, route)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("route-updated:%s", route.ID), route.ToJSON())
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) DeleteRoute(ctx context.Context, routeID string) error {
	// TODO: require auth

	route, err := h.DB.GetRoute(ctx, routeID)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf("route-deleted:%s", route.ID), route.ToJSON())
	if err != nil {
		return err
	}

	err = h.DB.DeleteRoute(ctx, routeID)
	if err != nil {
		return err
	}

	return nil
}

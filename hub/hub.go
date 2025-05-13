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
	// TODO: require auth
	err := h.DB.SetVehicle(ctx, vehicle)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemUpdatedTopic, "vehicle", vehicle.ID), vehicle.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range vehicle.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf(MapItemUpdatedTopic, "vehicle", cellID.ToToken()), vehicle.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) DeleteVehicle(ctx context.Context, vehicleID string) error {
	// TODO: require auth
	vehicle, err := h.DB.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	err = h.DB.DeleteVehicle(ctx, vehicle.ID)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemDeletedTopic, "vehicle", vehicle.ID), vehicle.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range vehicle.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf(MapItemDeletedTopic, "vehicle", cellID.ToToken()), vehicle.ToJSON())
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
	// TODO: require auth
	err := h.DB.SetStop(ctx, stop)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemUpdatedTopic, "stop", stop.ID), stop.ToJSON())
	if err != nil {
		return err
	}

	for _, cellID := range stop.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf(MapItemUpdatedTopic, "stop", cellID.ToToken()), stop.ToJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) DeleteStop(ctx context.Context, stopID string) error {
	// TODO: require auth
	stop, err := h.DB.GetStop(ctx, stopID)
	if err != nil {
		return err
	}

	err = h.DB.DeleteStop(ctx, stop.ID)
	if err != nil {
		return err
	}

	json := stop.ToJSON()
	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemDeletedTopic, "stop", stop.ID), json)
	if err != nil {
		return err
	}

	for _, cellID := range stop.Location.GetCellIDs() {
		err = h.PubSub.Publish(ctx, fmt.Sprintf(MapItemDeletedTopic, "stop", cellID.ToToken()), json)
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

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemUpdatedTopic, "trip", trip.ID), trip.ToJSON())
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) DeleteTrip(ctx context.Context, trip *models.Trip) error {
	// TODO: require auth
	err := h.DB.DeleteTrip(ctx, trip.ID)
	if err != nil {
		return err
	}

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemDeletedTopic, "trip", trip.ID), trip.ToJSON())
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

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemUpdatedTopic, "route", route.ID), route.ToJSON())
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

	err = h.PubSub.Publish(ctx, fmt.Sprintf(ItemDeletedTopic, "route", route.ID), route.ToJSON())
	if err != nil {
		return err
	}

	err = h.DB.DeleteRoute(ctx, routeID)
	if err != nil {
		return err
	}

	return nil
}

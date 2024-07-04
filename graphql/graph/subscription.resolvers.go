package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kiel-live/kiel-live/shared/models"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

// MapStopUpdated is the resolver for the mapStopUpdated field.
func (r *subscriptionResolver) MapStopUpdated(ctx context.Context, minLat float64, minLng float64, maxLat float64, maxLng float64) (<-chan *models.Stop, error) {
	ch := make(chan *models.Stop)

	r.subscribeBoundingBox(ctx, minLat, minLng, maxLat, maxLng, "map-stop-updated", func(message pubsub.Message) {
		stop := &models.Stop{}
		err := json.Unmarshal(message, stop)
		if err != nil {
			return
		}
		ch <- stop
	})

	return ch, nil
}

// MapStopDeleted is the resolver for the mapStopDeleted field.
func (r *subscriptionResolver) MapStopDeleted(ctx context.Context, minLat float64, minLng float64, maxLat float64, maxLng float64) (<-chan *models.Stop, error) {
	ch := make(chan *models.Stop)

	r.subscribeBoundingBox(ctx, minLat, minLng, maxLat, maxLng, "map-stop-deleted", func(message pubsub.Message) {
		stop := &models.Stop{}
		err := json.Unmarshal(message, stop)
		if err != nil {
			return
		}
		ch <- stop
	})

	return ch, nil
}

// MapVehicleUpdated is the resolver for the mapVehicleUpdated field.
func (r *subscriptionResolver) MapVehicleUpdated(ctx context.Context, minLat float64, minLng float64, maxLat float64, maxLng float64) (<-chan *models.Vehicle, error) {
	ch := make(chan *models.Vehicle)

	r.subscribeBoundingBox(ctx, minLat, minLng, maxLat, maxLng, "map-vehicle-updated", func(message pubsub.Message) {
		vehicle := &models.Vehicle{}
		err := json.Unmarshal(message, vehicle)
		if err != nil {
			return
		}
		ch <- vehicle
	})

	return ch, nil
}

// MapVehicleDeleted is the resolver for the mapVehicleDeleted field.
func (r *subscriptionResolver) MapVehicleDeleted(ctx context.Context, minLat float64, minLng float64, maxLat float64, maxLng float64) (<-chan *models.Vehicle, error) {
	ch := make(chan *models.Vehicle)

	r.subscribeBoundingBox(ctx, minLat, minLng, maxLat, maxLng, "map-vehicle-deleted", func(message pubsub.Message) {
		vehicle := &models.Vehicle{}
		err := json.Unmarshal(message, vehicle)
		if err != nil {
			return
		}
		ch <- vehicle
	})

	return ch, nil
}

// StopUpdated is the resolver for the stopUpdated field.
func (r *subscriptionResolver) StopUpdated(ctx context.Context, id string) (<-chan *models.Stop, error) {
	ch := make(chan *models.Stop)

	err := r.Hub.PubSub.Subscribe(ctx, fmt.Sprintf("stop-updated:%s", id), func(message pubsub.Message) {
		stop := &models.Stop{}
		err := json.Unmarshal(message, stop)
		if err != nil {
			return
		}
		ch <- stop
	})
	if err != nil {
		return nil, err
	}

	return ch, err
}

// StopDeleted is the resolver for the stopDeleted field.
func (r *subscriptionResolver) StopDeleted(ctx context.Context, id string) (<-chan *models.Stop, error) {
	ch := make(chan *models.Stop)

	err := r.Hub.PubSub.Subscribe(ctx, fmt.Sprintf("stop-deleted:%s", id), func(message pubsub.Message) {
		stop := &models.Stop{}
		err := json.Unmarshal(message, stop)
		if err != nil {
			return
		}
		ch <- stop
	})

	if err != nil {
		return nil, err
	}

	return ch, err
}

// VehicleUpdated is the resolver for the vehicleUpdated field.
func (r *subscriptionResolver) VehicleUpdated(ctx context.Context, id string) (<-chan *models.Vehicle, error) {
	ch := make(chan *models.Vehicle)

	err := r.Hub.PubSub.Subscribe(ctx, fmt.Sprintf("vehicle-updated:%s", id), func(message pubsub.Message) {
		vehicle := &models.Vehicle{}
		err := json.Unmarshal(message, vehicle)
		if err != nil {
			return
		}
		ch <- vehicle
	})

	if err != nil {
		return nil, err
	}

	return ch, err
}

// VehicleDeleted is the resolver for the vehicleDeleted field.
func (r *subscriptionResolver) VehicleDeleted(ctx context.Context, id string) (<-chan *models.Vehicle, error) {
	ch := make(chan *models.Vehicle)

	err := r.Hub.PubSub.Subscribe(ctx, fmt.Sprintf("vehicle-deleted:%s", id), func(message pubsub.Message) {
		vehicle := &models.Vehicle{}
		err := json.Unmarshal(message, vehicle)
		if err != nil {
			return
		}
		ch <- vehicle
	})

	if err != nil {
		return nil, err
	}

	return ch, err
}

// RouteUpdated is the resolver for the routeUpdated field.
func (r *subscriptionResolver) RouteUpdated(ctx context.Context, id string) (<-chan *models.Route, error) {
	panic(fmt.Errorf("not implemented: RouteUpdated - routeUpdated"))
}

// RouteDeleted is the resolver for the routeDeleted field.
func (r *subscriptionResolver) RouteDeleted(ctx context.Context, id string) (<-chan *models.Route, error) {
	panic(fmt.Errorf("not implemented: RouteDeleted - routeDeleted"))
}

// TripUpdated is the resolver for the tripUpdated field.
func (r *subscriptionResolver) TripUpdated(ctx context.Context, id string) (<-chan *models.Trip, error) {
	panic(fmt.Errorf("not implemented: TripUpdated - tripUpdated"))
}

// TripDeleted is the resolver for the tripDeleted field.
func (r *subscriptionResolver) TripDeleted(ctx context.Context, id string) (<-chan *models.Trip, error) {
	panic(fmt.Errorf("not implemented: TripDeleted - tripDeleted"))
}

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
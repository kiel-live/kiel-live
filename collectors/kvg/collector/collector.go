package collector

import (
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
)

type Collector interface {
	Run()
}

func NewCollector(client *client.Client, collectorType string, subscriptions *subscriptions.Subscriptions) (Collector, error) {
	// 	if c.channelType == "stops" {
	// 		return api.GetStops()
	// 	}

	// 	if c.channelType == "stop" {
	// 		return api.GetStop(c.entityID)
	// 	}

	// 	if c.channelType == "vehicles" {
	// 		return api.GetVehicles()
	// 	}

	// 	if c.channelType == "vehicle" {
	// 		return api.GetVehicle(c.entityID)
	// 	}

	// 	if c.channelType == "trip" {
	// 		return api.GetTrip(c.entityID)
	// 	}

	// 	if c.channelType == "route" {
	// 		return api.GetRoute(c.entityID)
	// 	}

	switch collectorType {
	case "map-vehicles":
		return &VehicleCollector{
			client: client,
		}, nil
	case "map-stops":
		return &StopCollector{
			client:        client,
			subscriptions: subscriptions,
		}, nil
	case "trips":
		return &TripCollector{
			client:        client,
			subscriptions: subscriptions,
		}, nil
	}

	return nil, fmt.Errorf("Collector type not supported")
}

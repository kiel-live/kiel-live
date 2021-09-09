package api

import (
	"encoding/json"
	"log"

	"github.com/kiel-live/kiel-live/protocol"
	"github.com/thoas/go-funk"
)

type vehicle struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Heading   int    `json:"heading"`
	Latitude  int    `json:"latitude"`
	Longitude int    `json:"longitude"`
	TripID    string `json:"tripId"`
}

type vehicles struct {
	Vehicles []vehicle `json:"vehicles"`
}

func (v *vehicle) parse() protocol.Vehicle {
	return protocol.Vehicle{
		ID:       "kvg-" + v.ID,
		Provider: "kvg", // TODO
		Name:     v.Name,
		Type:     protocol.VehicleTypeBus,
		State:    "onfire", // TODO
		Location: protocol.Location{
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		},
	}
}

func GetVehicles() []protocol.Vehicle {
	body, _ := post(vehiclesURL, nil)
	var vehicles vehicles
	if err := json.Unmarshal(body, &vehicles); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}

	// filter in-active vehicles
	vehicles.Vehicles = funk.Filter(vehicles.Vehicles, func(vehicle vehicle) bool {
		return vehicle.Latitude != 0
	}).([]vehicle)

	var res []protocol.Vehicle

	for _, vehicle := range vehicles.Vehicles {
		res = append(res, vehicle.parse())
	}

	return res
}

func GetVehicle(vehicleID string) protocol.Vehicle {
	// TODO
	return protocol.Vehicle{}
}

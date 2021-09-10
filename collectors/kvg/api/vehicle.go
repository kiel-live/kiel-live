package api

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	IsDeleted bool   `json:"isDeleted"`
}

type vehicles struct {
	Vehicles []vehicle `json:"vehicles"`
}

func (v *vehicle) parse() protocol.Vehicle {
	return protocol.Vehicle{
		ID:       "kvg" + v.ID,
		Provider: "kvg",
		Name:     v.Name,
		Type:     protocol.VehicleTypeBus,
		State:    "onfire", // TODO
		Location: protocol.Location{
			Heading:   v.Heading,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		},
	}
}

func GetVehicles() (res map[string]*protocol.Vehicle) {
	res = make(map[string]*protocol.Vehicle)
	url := fmt.Sprintf("%s?cacheBuster=%d&positionType=RAW", vehiclesURL, time.Now().Unix())
	body, _ := post(url, nil)
	// 	cacheBuster: new Date().getTime(),
	//   colorType: 'ROUTE_BASED',
	//   // lastUpdate: new Date().getTime(),
	//   positionType: 'RAW',
	var vehicles vehicles
	if err := json.Unmarshal(body, &vehicles); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}

	// filter in-active vehicles
	vehicles.Vehicles = funk.Filter(vehicles.Vehicles, func(vehicle vehicle) bool {
		return !vehicle.IsDeleted || vehicle.Latitude != 0
	}).([]vehicle)

	for _, vehicle := range vehicles.Vehicles {
		v := vehicle.parse()
		res[v.ID] = &v
	}

	return res
}

func GetVehicle(vehicleID string) protocol.Vehicle {
	// TODO
	return protocol.Vehicle{}
}

package api

import (
	"encoding/json"
	"fmt"
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
		ID:       IDPrefix + v.ID,
		Provider: "kvg",
		Name:     v.Name,
		Type:     protocol.VehicleTypeBus,
		TripID:   IDPrefix + v.TripID,
		State:    "onfire", // TODO
		Location: protocol.Location{
			Heading:   v.Heading,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		},
	}
}

func GetVehicles() (res map[string]*protocol.Vehicle, err error) {
	res = make(map[string]*protocol.Vehicle)
	url := fmt.Sprintf("%s?cacheBuster=%d&positionType=RAW", vehiclesURL, time.Now().Unix())
	body, err := post(url, nil)
	if err != nil {
		return nil, err
	}
	// 	cacheBuster: new Date().getTime(),
	//   colorType: 'ROUTE_BASED',
	//   // lastUpdate: new Date().getTime(),
	//   positionType: 'RAW',
	var vehicles vehicles
	if err := json.Unmarshal(body, &vehicles); err != nil {
		return nil, err
	}

	// filter in-active vehicles
	vehicles.Vehicles = funk.Filter(vehicles.Vehicles, func(vehicle vehicle) bool {
		return !vehicle.IsDeleted || vehicle.Latitude != 0
	}).([]vehicle)

	for _, vehicle := range vehicles.Vehicles {
		v := vehicle.parse()
		res[v.ID] = &v
	}

	return res, nil
}

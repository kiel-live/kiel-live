package collector

import (
	"encoding/json"
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/protocol"
)

type VehicleCollector struct {
	client *client.Client
}

func (c *VehicleCollector) publish(vehicle protocol.Vehicle) error {
	subject := fmt.Sprintf(protocol.SubjectMapVehicle, vehicle.Location.Longitude, vehicle.Location.Latitude, vehicle.ID)

	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		return err
	}

	err = c.client.Publish(subject, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (c *VehicleCollector) Run() {
	vehicles := api.GetVehicles()
	for _, vehicle := range vehicles {
		c.publish(vehicle)
	}
}

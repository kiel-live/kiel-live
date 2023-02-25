package collector

import (
	"encoding/json"
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/sirupsen/logrus"
)

type VehicleCollector struct {
	client   *client.Client
	vehicles map[string]*protocol.Vehicle
}

func isSameLocation(a protocol.Location, b protocol.Location) bool {
	return a.Heading == b.Heading && a.Latitude == b.Latitude && a.Longitude == b.Longitude
}

func isSameVehicle(a *protocol.Vehicle, b *protocol.Vehicle) bool {
	return a != nil && b != nil &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.ID == b.ID &&
		isSameLocation(a.Location, b.Location) &&
		a.Battery == b.Battery &&
		a.State == b.State &&
		a.TripID == b.TripID &&
		a.Type == b.Type
}

// returns list of changed or newly added vehicles
func (c *VehicleCollector) getChangedVehicles(vehicles map[string]*protocol.Vehicle) (changed []*protocol.Vehicle) {
	for _, v := range vehicles {
		if !isSameVehicle(v, c.vehicles[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *VehicleCollector) getRemovedVehicles(vehicles map[string]*protocol.Vehicle) (removed []*protocol.Vehicle) {
	for _, v := range c.vehicles {
		if _, ok := vehicles[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *VehicleCollector) publish(vehicle *protocol.Vehicle) error {
	subject := fmt.Sprintf(protocol.SubjectMapVehicle, vehicle.ID)

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

func (c *VehicleCollector) publishRemoved(vehicle *protocol.Vehicle) error {
	subject := fmt.Sprintf(protocol.SubjectMapVehicle, vehicle.ID)

	err := c.client.Publish(subject, string(protocol.DeletePayload))
	if err != nil {
		return err
	}

	return nil
}

func (c *VehicleCollector) SubjectsToIDs(subjects []string) []string {
	return []string{}
}

func (c *VehicleCollector) Run(_ []string, _ bool) {
	log := logrus.WithField("collector", "vehicle")
	vehicles, err := api.GetVehicles()
	if err != nil {
		log.Error(err)
		return
	}

	// publish all changed vehicles
	changed := c.getChangedVehicles(vehicles)
	for _, vehicle := range changed {
		err := c.publish(vehicle)
		if err != nil {
			log.Error(err)
		}
	}

	// publish all removed vehicles
	removed := c.getRemovedVehicles(vehicles)
	for _, vehicle := range removed {
		err := c.publishRemoved(vehicle)
		if err != nil {
			log.Error(err)
		}
	}

	log.Debugf("changed %d vehicles and removed %d", len(changed), len(removed))

	// update list of vehicles
	c.vehicles = vehicles
}

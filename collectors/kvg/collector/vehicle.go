package collector

import (
	"log/slog"
	"sync"

	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

type VehicleCollector struct {
	client   client.Client
	vehicles map[string]*models.Vehicle
	sync.Mutex
}

func isSameLocation(a *models.Location, b *models.Location) bool {
	return a.Heading == b.Heading && a.Latitude == b.Latitude && a.Longitude == b.Longitude
}

func isSameVehicle(a *models.Vehicle, b *models.Vehicle) bool {
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
func (c *VehicleCollector) getChangedVehicles(vehicles map[string]*models.Vehicle) (changed []*models.Vehicle) {
	for _, v := range vehicles {
		if !isSameVehicle(v, c.vehicles[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *VehicleCollector) getRemovedVehicles(vehicles map[string]*models.Vehicle) (removed []*models.Vehicle) {
	for _, v := range c.vehicles {
		if _, ok := vehicles[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *VehicleCollector) TopicToID(string) string {
	return ""
}

func (c *VehicleCollector) Run() {
	c.Lock()
	defer c.Unlock()

	log := slog.With("collector", "vehicle")
	vehicles, err := api.GetVehicles()
	if err != nil {
		log.Error("failed to get vehicles", "error", err)
		return
	}

	// publish all changed vehicles
	changed := c.getChangedVehicles(vehicles)
	for _, vehicle := range changed {
		log.Debug("publish updated vehicle", "vehicle", vehicle)
		err := c.client.UpdateVehicle(vehicle)
		if err != nil {
			log.Error("failed to update vehicle", "error", err)
		}
	}

	// publish all removed vehicles
	removed := c.getRemovedVehicles(vehicles)
	for _, vehicle := range removed {
		log.Debug("publish removed vehicle", "vehicle", vehicle)
		err := c.client.DeleteVehicle(vehicle.ID)
		if err != nil {
			log.Error("failed to delete vehicle", "error", err)
		}
	}

	log.Debug("collector run complete", "changed", len(changed), "removed", len(removed))

	// update list of vehicles
	c.vehicles = vehicles
}

func (c *VehicleCollector) RunSingle(string) {}

func (c *VehicleCollector) Reset() {
	c.Lock()
	defer c.Unlock()

	c.vehicles = make(map[string]*models.Vehicle)
}

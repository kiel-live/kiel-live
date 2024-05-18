package database

import (
	"errors"
	"sync"

	"github.com/kiel-live/kiel-live/hub/graph/model"

	"github.com/golang/geo/s2"
)

type MemoryDatabase struct {
	sync.RWMutex

	stops           map[string]*model.Stop
	stopsCellsIndex map[s2.CellID]map[string]struct{}

	vehicles           map[string]*model.Vehicle
	vehiclesCellsIndex map[s2.CellID]map[string]struct{}
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{
		stops:              make(map[string]*model.Stop),
		stopsCellsIndex:    make(map[s2.CellID]map[string]struct{}),
		vehicles:           make(map[string]*model.Vehicle),
		vehiclesCellsIndex: make(map[s2.CellID]map[string]struct{}),
	}
}

func (b *MemoryDatabase) Open() error {
	return nil
}

func (b *MemoryDatabase) Close() error {
	return nil
}

func (b *MemoryDatabase) GetStops(opts *ListOptions) ([]*model.Stop, error) {
	b.RLock()
	defer b.RUnlock()

	var stops []*model.Stop
	for _, cellID := range opts.Location.GetCellIDs() {
		if stopsInCell, ok := b.stopsCellsIndex[cellID]; ok {
			for stopID := range stopsInCell {
				if stop, ok := b.stops[stopID]; ok {
					stops = append(stops, stop)
				}
			}
		}
	}

	return stops, nil
}

func (b *MemoryDatabase) GetStop(id string) (*model.Stop, error) {
	b.RLock()
	defer b.RUnlock()

	if stop, ok := b.stops[id]; ok {
		return stop, nil
	}

	return nil, errors.New("stop not found")
}

func (b *MemoryDatabase) SetStop(stop *model.Stop) error {
	b.Lock()
	defer b.Unlock()

	cid := stop.Location.GetCellID()

	if oldStop, ok := b.stops[stop.ID]; ok {
		oldCid := oldStop.Location.GetCellID()
		if oldCid != cid {
			delete(b.stopsCellsIndex[oldCid], stop.ID)
			if len(b.stopsCellsIndex[oldCid]) == 0 {
				delete(b.stopsCellsIndex, oldCid)
			}
		}
	}

	b.stops[stop.ID] = stop

	if _, ok := b.stopsCellsIndex[cid]; !ok {
		b.stopsCellsIndex[cid] = make(map[string]struct{})
	}
	b.stopsCellsIndex[cid][stop.ID] = struct{}{}

	return nil
}

func (b *MemoryDatabase) DeleteStop(id string) error {
	b.Lock()
	defer b.Unlock()

	if stop, ok := b.stops[id]; ok {
		cid := stop.Location.GetCellID()
		delete(b.stopsCellsIndex[cid], id)
		if len(b.stopsCellsIndex[cid]) == 0 {
			delete(b.stopsCellsIndex, cid)
		}
		delete(b.stops, id)
	}

	return nil
}

func (b *MemoryDatabase) GetVehicles(opts *ListOptions) ([]*model.Vehicle, error) {
	b.RLock()
	defer b.RUnlock()

	var vehicles []*model.Vehicle
	for _, cellID := range opts.Location.GetCellIDs() {
		if vehiclesInCell, ok := b.vehiclesCellsIndex[cellID]; ok {
			for vehicleID := range vehiclesInCell {
				if vehicle, ok := b.vehicles[vehicleID]; ok {
					vehicles = append(vehicles, vehicle)
				}
			}
		}
	}

	return vehicles, nil
}

func (b *MemoryDatabase) GetVehicle(id string) (*model.Vehicle, error) {
	b.RLock()
	defer b.RUnlock()

	if vehicle, ok := b.vehicles[id]; ok {
		return vehicle, nil
	}

	return nil, errors.New("vehicle not found")
}

func (b *MemoryDatabase) SetVehicle(vehicle *model.Vehicle) error {
	b.Lock()
	defer b.Unlock()

	cid := vehicle.Location.GetCellID()

	if oldVehicle, ok := b.stops[vehicle.ID]; ok {
		oldCid := oldVehicle.Location.GetCellID()
		if oldCid != cid {
			delete(b.stopsCellsIndex[oldCid], vehicle.ID)
			if len(b.stopsCellsIndex[oldCid]) == 0 {
				delete(b.stopsCellsIndex, oldCid)
			}
		}
	}

	b.vehicles[vehicle.ID] = vehicle

	if _, ok := b.vehiclesCellsIndex[cid]; !ok {
		b.vehiclesCellsIndex[cid] = make(map[string]struct{})
	}
	b.vehiclesCellsIndex[cid][vehicle.ID] = struct{}{}

	return nil
}

func (b *MemoryDatabase) DeleteVehicle(id string) error {
	b.Lock()
	defer b.Unlock()

	if vehicle, ok := b.vehicles[id]; ok {
		cid := vehicle.Location.GetCellID()
		delete(b.vehiclesCellsIndex[cid], id)
		if len(b.vehiclesCellsIndex[cid]) == 0 {
			delete(b.vehiclesCellsIndex, cid)
		}
		delete(b.vehicles, id)
	}

	return nil
}

package database

import (
	"context"
	"sync"

	"github.com/kiel-live/kiel-live/pkg/models"

	"github.com/golang/geo/s2"
)

type MemoryDatabase struct {
	sync.RWMutex

	stops           map[string]*models.Stop
	stopsCellsIndex *CellIndex

	vehicles           map[string]*models.Vehicle
	vehiclesCellsIndex *CellIndex

	trips  map[string]*models.Trip
	routes map[string]*models.Route
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{
		stops:              make(map[string]*models.Stop),
		stopsCellsIndex:    NewCellIndex(),
		vehicles:           make(map[string]*models.Vehicle),
		vehiclesCellsIndex: NewCellIndex(),
		trips:              make(map[string]*models.Trip),
	}
}

func (b *MemoryDatabase) Open() error {
	return nil
}

func (b *MemoryDatabase) Close() error {
	return nil
}

func (b *MemoryDatabase) GetStops(_ context.Context, opts *ListOptions) ([]*models.Stop, error) {
	b.RLock()
	defer b.RUnlock()

	var stops []*models.Stop
	for _, cellID := range opts.Bounds.GetCellIDs() {
		stopIDs := b.stopsCellsIndex.GetItemIDs(cellID)
		for _, stopID := range stopIDs {
			if stop, ok := b.stops[stopID]; ok {
				stops = append(stops, stop)
			}
		}
	}

	return stops, nil
}

func (b *MemoryDatabase) GetStop(_ context.Context, id string) (*models.Stop, error) {
	b.RLock()
	defer b.RUnlock()

	if stop, ok := b.stops[id]; ok {
		return stop, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetStop(_ context.Context, stop *models.Stop) error {
	b.Lock()
	defer b.Unlock()

	cellID := stop.Location.GetCellID()

	var oldCellID s2.CellID
	if oldStop, ok := b.stops[stop.ID]; ok {
		oldCellID = oldStop.Location.GetCellID()
	}

	b.stops[stop.ID] = stop

	b.stopsCellsIndex.UpdateItem(stop.ID, []s2.CellID{cellID}, []s2.CellID{oldCellID})

	return nil
}

func (b *MemoryDatabase) DeleteStop(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if stop, ok := b.stops[id]; ok {
		b.stopsCellsIndex.RemoveItem(id, []s2.CellID{stop.Location.GetCellID()})
		delete(b.stops, id)
	}

	return nil
}

func (b *MemoryDatabase) GetVehicles(_ context.Context, opts *ListOptions) ([]*models.Vehicle, error) {
	b.RLock()
	defer b.RUnlock()

	var vehicles []*models.Vehicle
	for _, cellID := range opts.Bounds.GetCellIDs() {
		vehicleIDs := b.vehiclesCellsIndex.GetItemIDs(cellID)
		for _, vehicleID := range vehicleIDs {
			if vehicle, ok := b.vehicles[vehicleID]; ok {
				vehicles = append(vehicles, vehicle)
			}
		}
	}

	return vehicles, nil
}

func (b *MemoryDatabase) GetVehicle(_ context.Context, id string) (*models.Vehicle, error) {
	b.RLock()
	defer b.RUnlock()

	if vehicle, ok := b.vehicles[id]; ok {
		return vehicle, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetVehicle(_ context.Context, vehicle *models.Vehicle) error {
	b.Lock()
	defer b.Unlock()

	cellID := vehicle.Location.GetCellID()

	var oldCellID s2.CellID
	if oldVehicle, ok := b.vehicles[vehicle.ID]; ok {
		oldCellID = oldVehicle.Location.GetCellID()
	}

	b.vehicles[vehicle.ID] = vehicle

	b.vehiclesCellsIndex.UpdateItem(vehicle.ID, []s2.CellID{cellID}, []s2.CellID{oldCellID})

	return nil
}

func (b *MemoryDatabase) DeleteVehicle(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if vehicle, ok := b.vehicles[id]; ok {
		b.vehiclesCellsIndex.RemoveItem(id, []s2.CellID{vehicle.Location.GetCellID()})
		delete(b.vehicles, id)
	}

	return nil
}

func (b *MemoryDatabase) GetTrip(_ context.Context, id string) (*models.Trip, error) {
	b.RLock()
	defer b.RUnlock()

	if trip, ok := b.trips[id]; ok {
		return trip, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetTrip(_ context.Context, trip *models.Trip) error {
	b.Lock()
	defer b.Unlock()

	b.trips[trip.ID] = trip

	return nil
}

func (b *MemoryDatabase) DeleteTrip(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if trip, ok := b.trips[id]; ok {
		delete(b.trips, trip.ID)
	}

	return nil
}

func (b *MemoryDatabase) GetRoute(_ context.Context, id string) (*models.Route, error) {
	b.RLock()
	defer b.RUnlock()

	if route, ok := b.routes[id]; ok {
		return route, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetRoute(_ context.Context, route *models.Route) error {
	b.Lock()
	defer b.Unlock()

	b.routes[route.ID] = route

	return nil
}

func (b *MemoryDatabase) DeleteRoute(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if route, ok := b.routes[id]; ok {
		delete(b.routes, route.ID)
	}

	return nil
}

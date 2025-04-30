package database

import (
	"context"
	"errors"
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
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{
		stops:              make(map[string]*models.Stop),
		stopsCellsIndex:    NewCellIndex(),
		vehicles:           make(map[string]*models.Vehicle),
		vehiclesCellsIndex: NewCellIndex(),
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
	for _, cellID := range opts.Location.GetCellIDs() {
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

	return nil, errors.New("stop not found")
}

func (b *MemoryDatabase) SetStop(_ context.Context, stop *models.Stop) error {
	b.Lock()
	defer b.Unlock()

	cids := stop.Location.GetCellIDs()

	oldCids := []s2.CellID{}
	if oldStop, ok := b.stops[stop.ID]; ok {
		oldCids = oldStop.Location.GetCellIDs()
	}

	b.stops[stop.ID] = stop

	b.stopsCellsIndex.UpdateItem(stop.ID, cids, oldCids)

	return nil
}

func (b *MemoryDatabase) DeleteStop(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if stop, ok := b.stops[id]; ok {
		b.stopsCellsIndex.RemoveItem(id, stop.Location.GetCellIDs())
		delete(b.stops, id)
	}

	return nil
}

func (b *MemoryDatabase) GetVehicles(_ context.Context, opts *ListOptions) ([]*models.Vehicle, error) {
	b.RLock()
	defer b.RUnlock()

	var vehicles []*models.Vehicle
	for _, cellID := range opts.Location.GetCellIDs() {
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

	return nil, errors.New("vehicle not found")
}

func (b *MemoryDatabase) SetVehicle(_ context.Context, vehicle *models.Vehicle) error {
	b.Lock()
	defer b.Unlock()

	cids := vehicle.Location.GetCellIDs()

	oldCids := []s2.CellID{}
	if oldVehicle, ok := b.stops[vehicle.ID]; ok {
		oldCids = oldVehicle.Location.GetCellIDs()
	}

	b.vehicles[vehicle.ID] = vehicle

	b.vehiclesCellsIndex.UpdateItem(vehicle.ID, cids, oldCids)

	return nil
}

func (b *MemoryDatabase) DeleteVehicle(_ context.Context, id string) error {
	b.Lock()
	defer b.Unlock()

	if vehicle, ok := b.vehicles[id]; ok {
		b.vehiclesCellsIndex.RemoveItem(id, vehicle.Location.GetCellIDs())
		delete(b.vehicles, id)
	}

	return nil
}

type CellIndex struct {
	sync.RWMutex
	index map[s2.CellID]map[string]struct{}
}

func NewCellIndex() *CellIndex {
	return &CellIndex{
		index: make(map[s2.CellID]map[string]struct{}),
	}
}

func (c *CellIndex) UpdateItem(itemID string, newIDs []s2.CellID, oldIDs []s2.CellID) {
	c.Lock()
	defer c.Unlock()

	toDelete := make(map[s2.CellID]struct{})
	for _, oldID := range oldIDs {
		toDelete[oldID] = struct{}{}
	}

	toAdd := make(map[s2.CellID]struct{})
	for _, newID := range newIDs {
		if _, ok := toDelete[newID]; ok {
			// keep item and therefore don't delete it
			delete(toDelete, newID)
		} else {
			toAdd[newID] = struct{}{}
		}
	}

	for id := range toDelete {
		if _, ok := c.index[id]; ok {
			delete(c.index[id], itemID)

			// delete empty cell
			if len(c.index[id]) == 0 {
				delete(c.index, id)
			}
		}
	}

	for id := range toAdd {
		// create new cell if it doesn't exist
		if _, ok := c.index[id]; !ok {
			c.index[id] = make(map[string]struct{})
		}
		c.index[id][itemID] = struct{}{}
	}
}

func (c *CellIndex) AddItem(itemID string, cellIDs []s2.CellID) {
	c.UpdateItem(itemID, cellIDs, nil)
}

func (c *CellIndex) RemoveItem(itemID string, cellIDs []s2.CellID) {
	c.UpdateItem(itemID, nil, cellIDs)
}

func (c *CellIndex) GetItemIDs(cellID s2.CellID) []string {
	c.RLock()
	defer c.RUnlock()

	itemIDs, ok := c.index[cellID]
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(itemIDs))
	for id := range itemIDs {
		ids = append(ids, id)
	}

	return ids
}

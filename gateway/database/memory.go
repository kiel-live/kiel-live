package database

import (
	"context"

	"github.com/kiel-live/kiel-live/pkg/database/memory"
	"github.com/kiel-live/kiel-live/pkg/models"

	"github.com/golang/geo/s2"
)

type MemoryDatabase struct {
	vehicles *memory.SpatialStore[*models.Vehicle]
	stops    *memory.SpatialStore[*models.Stop]
	trips    *memory.Store[*models.Trip]
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{
		vehicles: memory.NewSpatialStore(
			func(v *models.Vehicle) string { return v.ID },
			func(v *models.Vehicle) s2.CellID { return v.Location.GetCellID() },
		),
		stops: memory.NewSpatialStore(
			func(s *models.Stop) string { return s.ID },
			func(s *models.Stop) s2.CellID { return s.Location.GetCellID() },
		),
		trips: memory.NewStore(
			func(t *models.Trip) string { return t.ID },
		),
	}
}

func (b *MemoryDatabase) Open() error {
	return nil
}

func (b *MemoryDatabase) Close() error {
	return nil
}

func (b *MemoryDatabase) GetStops(_ context.Context, opts *ListOptions) ([]*models.Stop, error) {
	return b.stops.GetInBounds(opts.Bounds.GetCellIDs()), nil
}

func (b *MemoryDatabase) GetStop(_ context.Context, id string) (*models.Stop, error) {
	if stop, ok := b.stops.Get(id); ok {
		return stop, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetStop(_ context.Context, stop *models.Stop) error {
	b.stops.Set(stop)
	return nil
}

func (b *MemoryDatabase) DeleteStop(_ context.Context, id string) error {
	if _, ok := b.stops.Get(id); ok {
		b.stops.Delete(id)
	}

	return nil
}

func (b *MemoryDatabase) GetVehicles(_ context.Context, opts *ListOptions) ([]*models.Vehicle, error) {
	return b.vehicles.GetInBounds(opts.Bounds.GetCellIDs()), nil
}

func (b *MemoryDatabase) GetVehicle(_ context.Context, id string) (*models.Vehicle, error) {
	if vehicle, ok := b.vehicles.Get(id); ok {
		return vehicle, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetVehicle(_ context.Context, vehicle *models.Vehicle) error {
	b.vehicles.Set(vehicle)
	return nil
}

func (b *MemoryDatabase) DeleteVehicle(_ context.Context, id string) error {
	if _, ok := b.vehicles.Get(id); ok {
		b.vehicles.Delete(id)
	}

	return nil
}

func (b *MemoryDatabase) GetTrip(_ context.Context, id string) (*models.Trip, error) {
	if trip, ok := b.trips.Get(id); ok {
		return trip, nil
	}

	return nil, ErrItemNotFound
}

func (b *MemoryDatabase) SetTrip(_ context.Context, trip *models.Trip) error {
	b.trips.Set(trip)

	return nil
}

func (b *MemoryDatabase) DeleteTrip(_ context.Context, id string) error {
	b.trips.Delete(id)

	return nil
}

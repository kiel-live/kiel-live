package database

import (
	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/hub/graph/model"
)

type BoundingBox struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
}

func (b *BoundingBox) GetCellIDs() []s2.CellID {
	p1 := s2.LatLngFromDegrees(b.MinLat, b.MinLng)
	p2 := s2.LatLngFromDegrees(b.MaxLat, b.MaxLng)
	r := s2.RectFromLatLng(p1).AddPoint(p2)
	rc := s2.RegionCoverer{MaxLevel: 11, MinLevel: 11, MaxCells: 100}
	return rc.Covering(r)
}

type ListOptions struct {
	Location *BoundingBox
	Limit    int
}

type Database interface {
	Open() error
	Close() error

	// Stops
	GetStops(*ListOptions) ([]*model.Stop, error)
	GetStop(id string) (*model.Stop, error)
	SetStop(stop *model.Stop) error
	DeleteStop(id string) error

	// Vehicles
	GetVehicles(*ListOptions) ([]*model.Vehicle, error)
	GetVehicle(id string) (*model.Vehicle, error)
	SetVehicle(vehicle *model.Vehicle) error
	DeleteVehicle(id string) error
}

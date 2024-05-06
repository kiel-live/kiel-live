package database

import "github.com/kiel-live/kiel-live/hub/graph/model"

type BoundingBox struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
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

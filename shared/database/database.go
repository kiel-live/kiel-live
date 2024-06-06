package database

import (
	"github.com/kiel-live/kiel-live/shared/models"
)

type ListOptions struct {
	Location *models.BoundingBox
	Limit    int
}

type Database interface {
	Open() error
	Close() error

	// Stops
	GetStops(*ListOptions) ([]*models.Stop, error)
	GetStop(id string) (*models.Stop, error)
	SetStop(stop *models.Stop) error
	DeleteStop(id string) error

	// Vehicles
	GetVehicles(*ListOptions) ([]*models.Vehicle, error)
	GetVehicle(id string) (*models.Vehicle, error)
	SetVehicle(vehicle *models.Vehicle) error
	DeleteVehicle(id string) error
}

package database

import (
	"context"

	"github.com/kiel-live/kiel-live/pkg/models"
)

type ListOptions struct {
	Location *models.BoundingBox
	Limit    int
}

type Database interface {
	Open() error
	Close() error

	// Stops
	GetStops(ctx context.Context, opts *ListOptions) ([]*models.Stop, error)
	GetStop(ctx context.Context, id string) (*models.Stop, error)
	SetStop(ctx context.Context, stop *models.Stop) error
	DeleteStop(ctx context.Context, id string) error

	// Vehicles
	GetVehicles(ctx context.Context, opts *ListOptions) ([]*models.Vehicle, error)
	GetVehicle(ctx context.Context, id string) (*models.Vehicle, error)
	SetVehicle(ctx context.Context, vehicle *models.Vehicle) error
	DeleteVehicle(ctx context.Context, id string) error
}

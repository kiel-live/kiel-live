package search

import (
	"context"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/search"
)

type Search interface {
	// Stops
	SetStop(ctx context.Context, stop *models.Stop) error
	DeleteStop(ctx context.Context, id string) error

	// Vehicles
	SetVehicle(ctx context.Context, vehicle *models.Vehicle) error
	DeleteVehicle(ctx context.Context, id string) error

	// Search
	Search(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
}

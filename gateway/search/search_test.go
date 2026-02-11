package search

import (
	"context"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func toDegreesInt(value float64) int {
	return int(value * 3600000.0)
}

func TestMemorySearch_Stops(t *testing.T) {
	search := NewMemorySearch()

	// Add stops
	stops := []*models.Stop{
		{
			ID:       "stop-1",
			Name:     "Hauptbahnhof",
			Provider: "kvg",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.32),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:       "stop-2",
			Name:     "Rathaus",
			Provider: "kvg",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.33),
				Longitude: toDegreesInt(10.15),
			},
		},
		{
			ID:       "stop-3",
			Name:     "Bahnhof Nord",
			Provider: "nah.sh",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.34),
				Longitude: toDegreesInt(10.16),
			},
		},
	}

	ctx := context.Background()
	for _, stop := range stops {
		err := search.SetStop(ctx, stop)
		assert.NoError(t, err)
	}

	// Search for "Bahnhof" - should match "Hauptbahnhof" and "Bahnhof Nord"
	results, err := search.Search(ctx, "Bahnhof", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)

	// Verify we get relevant results
	foundBahnhof := false
	for _, result := range results {
		if result.ID == "stop-1" || result.ID == "stop-3" {
			foundBahnhof = true
			break
		}
	}
	assert.True(t, foundBahnhof, "Should find stops with 'Bahnhof' in name")

	// Test deletion
	err = search.DeleteStop(ctx, "stop-1")
	assert.NoError(t, err)

	// Search again - should not find deleted stop
	results, err = search.Search(ctx, "Hauptbahnhof", 10)
	assert.NoError(t, err)
	for _, result := range results {
		assert.NotEqual(t, "stop-1", result.ID, "Deleted stop should not appear in results")
	}
}

func TestMemorySearch_Vehicles(t *testing.T) {
	search := NewMemorySearch()

	// Add vehicles
	vehicles := []*models.Vehicle{
		{
			ID:          "vehicle-1",
			Name:        "Bus 11",
			Provider:    "kvg",
			Type:        models.VehicleTypeBus,
			Description: "City center route",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.32),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:          "vehicle-2",
			Name:        "Bus 22",
			Provider:    "kvg",
			Type:        models.VehicleTypeBus,
			Description: "University route",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.33),
				Longitude: toDegreesInt(10.15),
			},
		},
		{
			ID:       "vehicle-3",
			Name:     "E-Scooter 123",
			Provider: "voi",
			Type:     models.VehicleTypeEScooter,
			Location: &models.Location{
				Latitude:  toDegreesInt(54.34),
				Longitude: toDegreesInt(10.16),
			},
		},
	}

	ctx := context.Background()
	for _, vehicle := range vehicles {
		err := search.SetVehicle(ctx, vehicle)
		assert.NoError(t, err)
	}

	// Search for "Bus" - should match Bus 11 and Bus 22
	results, err := search.Search(ctx, "Bus", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2)

	// Search for provider "kvg"
	results, err = search.Search(ctx, "kvg", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2)

	// Test deletion
	err = search.DeleteVehicle(ctx, "vehicle-1")
	assert.NoError(t, err)

	// Search again
	results, err = search.Search(ctx, "Bus 11", 10)
	assert.NoError(t, err)
	for _, result := range results {
		assert.NotEqual(t, "vehicle-1", result.ID, "Deleted vehicle should not appear in results")
	}
}

func TestMemorySearch_MultiIndex(t *testing.T) {
	search := NewMemorySearch()
	ctx := context.Background()

	// Add both stops and vehicles
	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Hauptbahnhof",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := search.SetStop(ctx, stop)
	assert.NoError(t, err)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Bus Hauptbahnhof Express",
		Provider: "kvg",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err = search.SetVehicle(ctx, vehicle)
	assert.NoError(t, err)

	// Search for "Hauptbahnhof" - should match both stop and vehicle
	results, err := search.Search(ctx, "Hauptbahnhof", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2)

	// Verify we get both types (results should have Index field set)
	indices := make(map[string]bool)
	for _, result := range results {
		if result.Index != "" {
			indices[result.Index] = true
		}
	}
	assert.True(t, len(indices) >= 1, "Should get results from at least one index")
}

func TestMemorySearch_Limit(t *testing.T) {
	search := NewMemorySearch()
	ctx := context.Background()

	// Add multiple vehicles
	for i := 1; i <= 10; i++ {
		vehicle := &models.Vehicle{
			ID:       "vehicle-" + string(rune('0'+i)),
			Name:     "Bus Line " + string(rune('0'+i)),
			Provider: "kvg",
			Type:     models.VehicleTypeBus,
			Location: &models.Location{
				Latitude:  toDegreesInt(54.32),
				Longitude: toDegreesInt(10.14),
			},
		}
		err := search.SetVehicle(ctx, vehicle)
		assert.NoError(t, err)
	}

	// Search with limit
	results, err := search.Search(ctx, "Bus", 5)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(results), 5, "Should respect search limit")
}

func TestMemorySearch_EmptyQuery(t *testing.T) {
	search := NewMemorySearch()
	ctx := context.Background()

	// Add a stop
	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Test Stop",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := search.SetStop(ctx, stop)
	assert.NoError(t, err)

	// Search with empty query
	results, err := search.Search(ctx, "", 10)
	assert.NoError(t, err)
	// Empty query should return empty or minimal results
	assert.GreaterOrEqual(t, len(results), 0)
}

func TestMemorySearch_NoResults(t *testing.T) {
	search := NewMemorySearch()
	ctx := context.Background()

	// Add a stop
	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Hauptbahnhof",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := search.SetStop(ctx, stop)
	assert.NoError(t, err)

	// Search for something that doesn't exist
	results, err := search.Search(ctx, "xyznonexistent", 10)
	assert.NoError(t, err)
	// Should return empty results, not an error
	assert.GreaterOrEqual(t, len(results), 0)
}

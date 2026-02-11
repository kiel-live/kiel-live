package database

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func toDegreesInt(value float64) int {
	return int(value * 3600000.0)
}

func TestMemoryDatabase(t *testing.T) {
	db := NewMemoryDatabase()
	err := db.Open()
	assert.NoError(t, err)
	defer db.Close()

	err = db.SetStop(t.Context(), &models.Stop{
		ID: "1",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.31981897337084),
			Longitude: toDegreesInt(10.182968719044112),
		},
		Name: "Central Station",
	})
	assert.NoError(t, err)

	stop, err := db.GetStop(t.Context(), "1")
	assert.NoError(t, err)
	assert.Equal(t, "Central Station", stop.Name)

	stops, err := db.GetStops(t.Context(), &ListOptions{
		Bounds: &models.BoundingBox{
			North: 54.296181,
			East:  10.107290,
			South: 54.345022,
			West:  10.196574,
		},
	})
	assert.NoError(t, err)
	assert.Len(t, stops, 1)
	assert.Equal(t, "Central Station", stops[0].Name)

	err = db.DeleteStop(t.Context(), "1")
	assert.NoError(t, err)

	_, err = db.GetStop(t.Context(), "1")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)
}

func TestMemoryDatabase_Vehicles(t *testing.T) {
	db := NewMemoryDatabase()
	err := db.Open()
	assert.NoError(t, err)
	defer db.Close()

	// Test SetVehicle
	vehicle := &models.Vehicle{
		ID:       "bus-123",
		Provider: "kvg",
		Name:     "Bus 11",
		Type:     models.VehicleTypeBus,
		State:    "active",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.3233),
			Longitude: toDegreesInt(10.1394),
		},
	}
	err = db.SetVehicle(t.Context(), vehicle)
	assert.NoError(t, err)

	// Test GetVehicle
	retrieved, err := db.GetVehicle(t.Context(), "bus-123")
	assert.NoError(t, err)
	assert.Equal(t, "bus-123", retrieved.ID)
	assert.Equal(t, "kvg", retrieved.Provider)
	assert.Equal(t, "Bus 11", retrieved.Name)
	assert.Equal(t, models.VehicleTypeBus, retrieved.Type)

	// Test GetVehicles with bounds
	vehicles, err := db.GetVehicles(t.Context(), &ListOptions{
		Bounds: &models.BoundingBox{
			North: 54.35,
			East:  10.15,
			South: 54.30,
			West:  10.12,
		},
	})
	assert.NoError(t, err)
	assert.Len(t, vehicles, 1)
	assert.Equal(t, "bus-123", vehicles[0].ID)

	// Test vehicle not found
	_, err = db.GetVehicle(t.Context(), "non-existent")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)

	// Test DeleteVehicle
	err = db.DeleteVehicle(t.Context(), "bus-123")
	assert.NoError(t, err)

	_, err = db.GetVehicle(t.Context(), "bus-123")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)

	// Test delete non-existent vehicle (should not error)
	err = db.DeleteVehicle(t.Context(), "non-existent")
	assert.NoError(t, err)
}

func TestMemoryDatabase_Trips(t *testing.T) {
	db := NewMemoryDatabase()
	err := db.Open()
	assert.NoError(t, err)
	defer db.Close()

	// Test SetTrip
	trip := &models.Trip{
		ID:        "trip-456",
		Provider:  "kvg",
		VehicleID: "bus-123",
		Direction: "Innenstadt",
		Arrivals: []*models.TripArrival{
			{
				ID:      "stop-1",
				Name:    "Hauptbahnhof",
				State:   models.Planned,
				Planned: "2024-01-01T10:00:00Z",
			},
			{
				ID:      "stop-2",
				Name:    "Rathaus",
				State:   models.Predicted,
				Planned: "2024-01-01T10:05:00Z",
			},
		},
	}
	err = db.SetTrip(t.Context(), trip)
	assert.NoError(t, err)

	// Test GetTrip
	retrieved, err := db.GetTrip(t.Context(), "trip-456")
	assert.NoError(t, err)
	assert.Equal(t, "trip-456", retrieved.ID)
	assert.Equal(t, "kvg", retrieved.Provider)
	assert.Equal(t, "bus-123", retrieved.VehicleID)
	assert.Equal(t, "Innenstadt", retrieved.Direction)
	assert.Len(t, retrieved.Arrivals, 2)
	assert.Equal(t, "Hauptbahnhof", retrieved.Arrivals[0].Name)

	// Test trip not found
	_, err = db.GetTrip(t.Context(), "non-existent")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)

	// Test DeleteTrip
	err = db.DeleteTrip(t.Context(), "trip-456")
	assert.NoError(t, err)

	_, err = db.GetTrip(t.Context(), "trip-456")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)

	// Test delete non-existent trip (should not error)
	err = db.DeleteTrip(t.Context(), "non-existent")
	assert.NoError(t, err)
}

func TestMemoryDatabase_MultipleBounds(t *testing.T) {
	db := NewMemoryDatabase()
	err := db.Open()
	assert.NoError(t, err)
	defer db.Close()

	// Add stops at different locations
	stops := []*models.Stop{
		{
			ID:   "stop-north",
			Name: "North Stop",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.35),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:   "stop-south",
			Name: "South Stop",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.30),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:   "stop-center",
			Name: "Center Stop",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.325),
				Longitude: toDegreesInt(10.14),
			},
		},
	}

	for _, stop := range stops {
		err = db.SetStop(t.Context(), stop)
		assert.NoError(t, err)
	}

	// Query for northern area
	northStops, err := db.GetStops(t.Context(), &ListOptions{
		Bounds: &models.BoundingBox{
			North: 54.36,
			East:  10.15,
			South: 54.33,
			West:  10.13,
		},
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(northStops), 1)

	// Query for all
	allStops, err := db.GetStops(t.Context(), &ListOptions{
		Bounds: &models.BoundingBox{
			North: 54.36,
			East:  10.15,
			South: 54.29,
			West:  10.13,
		},
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(allStops), 3)
}

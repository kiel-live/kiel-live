package database

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
)

func TestMemoryDatabase(t *testing.T) {
	db := NewMemoryDatabase()
	if err := db.Open(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err := db.SetStop(t.Context(), &models.Stop{
		ID: "1",
		Location: &models.Location{
			Latitude:  int(54.31981897337084 * 360000),
			Longitude: int(10.182968719044112 * 360000),
		},
		Name: "Central Station",
	})
	if err != nil {
		t.Fatal(err)
	}

	stop, err := db.GetStop(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	if stop.Name != "Central Station" {
		t.Fatalf("expected Central Station, got %s", stop.Name)
	}

	stops, err := db.GetStops(t.Context(), &ListOptions{
		Bounds: &models.BoundingBox{
			North: 54.296181,
			East:  10.107290,
			South: 54.345022,
			West:  10.196574,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(stops) != 1 {
		t.Fatalf("expected 1 stop, got %d", len(stops))
	}

	if stops[0].Name != "Central Station" {
		t.Fatalf("expected Central Station, got %s", stops[0].Name)
	}

	err = db.DeleteStop(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetStop(t.Context(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
}

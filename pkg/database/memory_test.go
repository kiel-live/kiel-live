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

	err := db.SetStop(&models.Stop{
		ID: "1",
		Location: &models.Location{
			Latitude:  54.31981897337084,
			Longitude: 10.182968719044112,
			Heading:   nil,
		},
		Name: "Central Station",
	})
	if err != nil {
		t.Fatal(err)
	}

	stop, err := db.GetStop("1")
	if err != nil {
		t.Fatal(err)
	}

	if stop.Name != "Central Station" {
		t.Fatalf("expected Central Station, got %s", stop.Name)
	}

	stops, err := db.GetStops(&ListOptions{
		Location: &models.BoundingBox{
			MinLat: 54.526130648172995,
			MinLng: 9.876994965672509,
			MaxLat: 53.95617973610979,
			MaxLng: 10.709999024470449,
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

	err = db.DeleteStop("1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetStop("1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetCellIDs(t *testing.T) {
	ids := (&models.BoundingBox{
		MinLat: 54.526130648172995,
		MinLng: 9.876994965672509,
		MaxLat: 53.95617973610979,
		MaxLng: 10.709999024470449,
	}).GetCellIDs()

	poiID := (&models.Location{
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
		Heading:   nil,
	}).GetCellID()

	found := false
	for _, id := range ids {
		if id == poiID {
			found = true
			break
		}
	}

	if found == false {
		t.Fatalf("expected %d in %v", poiID, ids)
	}
}

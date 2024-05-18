package database

import (
	"testing"

	"github.com/kiel-live/kiel-live/hub/graph/model"
)

func TestMemoryDatabase(t *testing.T) {
	db := NewMemoryDatabase()
	if err := db.Open(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err := db.SetStop(&model.Stop{
		ID: "1",
		Location: &model.Location{
			Latitude:  54.31981897337084,
			Longitude: 10.182968719044112,
			Heading:   32,
		},
		Name: "Hamburg",
	})
	if err != nil {
		t.Fatal(err)
	}

	stop, err := db.GetStop("1")
	if err != nil {
		t.Fatal(err)
	}

	if stop.Name != "Hamburg" {
		t.Fatalf("expected Hamburg, got %s", stop.Name)
	}

	stops, err := db.GetStops(&ListOptions{
		Location: &BoundingBox{
			MinLat: 54.1634014386689,
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

	if stops[0].Name != "Hamburg" {
		t.Fatalf("expected Hamburg, got %s", stops[0].Name)
	}

	stops, err = db.GetStops(&ListOptions{
		Location: &BoundingBox{
			MinLat: 54.1634014386689,
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

	err = db.DeleteStop("1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCellIDs(t *testing.T) {
	db := MemoryDatabase{}
	ids := db.getCellIDsForListOptions(&ListOptions{
		Location: &BoundingBox{
			MinLat: 54.526130648172995,
			MinLng: 9.876994965672509,
			MaxLat: 53.95617973610979,
			MaxLng: 10.709999024470449,
		},
	})

	poiID := (&model.Location{
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
		Heading:   32,
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

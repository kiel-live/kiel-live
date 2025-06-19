package models_test

import (
	"slices"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestLocationGetCellIDs(t *testing.T) {
	l := &models.Location{
		Latitude:  54.31981897337084 * 3600000,
		Longitude: 10.182968719044112 * 3600000,
		Heading:   nil,
	}

	cells := l.GetCellIDs()
	assert.Len(t, cells, 1) // max = min => 1
}

func TestBoundingBoxGetCellIDs(t *testing.T) {
	b := &models.BoundingBox{
		North: 54.296181,
		East:  10.107290,
		South: 54.345022,
		West:  10.196574,
	}

	cells := b.GetCellIDs()
	assert.Len(t, cells, 10)

	for _, cell := range cells {
		assert.GreaterOrEqual(t, cell.Level(), models.MinLevel)
		assert.LessOrEqual(t, cell.Level(), models.MaxLevel)
	}

	for _, cell := range cells {
		assert.GreaterOrEqual(t, cell.LatLng().Lat.Degrees(), b.North)
		assert.LessOrEqual(t, cell.LatLng().Lat, b.South)
		assert.GreaterOrEqual(t, cell.LatLng().Lng, b.East)
		assert.LessOrEqual(t, cell.LatLng().Lng, b.West)
	}
}

func TestGetCellIDs(t *testing.T) {
	ids := (&models.BoundingBox{
		North: 54.526130648172995,
		East:  9.876994965672509,
		South: 53.95617973610979,
		West:  10.709999024470449,
	}).GetCellIDs()

	poiID := (&models.Location{
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
		Heading:   nil,
	}).GetCellID()

	found := slices.Contains(ids, poiID)
	if found == false {
		t.Fatalf("expected %d in %v", poiID, ids)
	}
}

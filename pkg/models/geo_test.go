package models_test

import (
	"slices"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func toDegreesInt(value float64) int {
	return int(value * 3600000.0)
}

func TestLocationGetCellID(t *testing.T) {
	l := &models.Location{
		Latitude:  toDegreesInt(54.31981897337084),
		Longitude: toDegreesInt(10.182968719044112),
	}

	cell := l.GetCellID()
	assert.Equal(t, cell.ToToken(), "47b257")
}

func TestBoundingBoxGetCellIDs(t *testing.T) {
	b := &models.BoundingBox{
		North: 54.296181,
		East:  10.107290,
		South: 54.345022,
		West:  10.196574,
	}

	cells := b.GetCellIDs()
	assert.Len(t, cells, 2)

	for _, cell := range cells {
		assert.GreaterOrEqual(t, cell.Level(), models.MinLevel)
		assert.LessOrEqual(t, cell.Level(), models.MaxLevel)
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
		Latitude:  toDegreesInt(54.31981897337084),
		Longitude: toDegreesInt(10.182968719044112),
		Heading:   0,
	}).GetCellID()

	found := slices.Contains(ids, poiID)
	if found == false {
		t.Fatalf("expected %d in %v", poiID, ids)
	}
}

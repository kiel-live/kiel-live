package models_test

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestLocationGetCellIDs(t *testing.T) {
	l := &models.Location{
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
		Heading:   nil,
	}

	cells := l.GetCellIDs()
	assert.Len(t, cells, models.MaxLevel)
}

func TestBoundingBoxGetCellIDs(t *testing.T) {
	b := &models.BoundingBox{
		MinLat: 54.296181,
		MinLng: 10.107290,
		MaxLat: 54.345022,
		MaxLng: 10.196574,
	}

	cells := b.GetCellIDs()
	assert.Len(t, cells, 10)

	for _, cell := range cells {
		assert.GreaterOrEqual(t, cell.Level(), models.MinLevel)
		assert.LessOrEqual(t, cell.Level(), models.MaxLevel)
	}

	for _, cell := range cells {
		assert.GreaterOrEqual(t, cell.LatLng().Lat.Degrees(), b.MinLat)
		assert.LessOrEqual(t, cell.LatLng().Lat, b.MaxLat)
		assert.GreaterOrEqual(t, cell.LatLng().Lng, b.MinLng)
		assert.LessOrEqual(t, cell.LatLng().Lng, b.MaxLng)
	}
}

package models_test

import (
	"testing"

	"github.com/kiel-live/kiel-live/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestXxx1(t *testing.T) {
	l := &models.Location{
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
		Heading:   nil,
	}

	cells := l.GetCellIDs()
	assert.Len(t, cells, models.MaxLevel)
}

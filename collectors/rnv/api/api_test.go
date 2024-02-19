package api_test

import (
	"context"
	"testing"

	"github.com/kiel-live/kiel-live/collectors/rnv/api"
)

func TestGetStops(t *testing.T) {
	ctx := context.Background()
	stops, err := api.GetStops(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(stops)
}

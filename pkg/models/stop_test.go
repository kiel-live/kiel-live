package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestStop(t *testing.T) {
	// v := []byte("{\"id\":\"hello\",\"provider\":\"kvg\",\"name\":\"123\",\"type\":\"bike\",\"routes\":null,\"alerts\":null,\"arrivals\":null,\"vehicles\":null,\"location\":{\"latitude\":12.3,\"longitude\":54.7,\"heading\":0}}")
	// v := []byte("{\"id\":\"hello\",\"provider\":\"kvg\",\"name\":\"123\",\"type\":\"bike\",\"location\":{\"latitude\":12.3,\"longitude\":54.7,\"heading\":0}}")

	o := &models.Stop{
		ID:       "hello",
		Provider: "kvg",
		Name:     "123",
		Type:     "bike",
		Location: &models.Location{
			Latitude:  toDegreesInt(12.3),
			Longitude: toDegreesInt(54.7),
			Heading:   0,
		},
	}

	s, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	fmt.Println(string(s))

	stop := &models.Stop{}

	err = json.Unmarshal(s, stop)
	assert.NoError(t, err)
}

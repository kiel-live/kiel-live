package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"public-transport-app/pkg/models"

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
			Latitude:  12.3,
			Longitude: 54.7,
			Heading:   nil,
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

package testing

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
	"github.com/stretchr/testify/assert"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type StopInput struct {
	ID       string   `json:"id"`
	Provider string   `json:"provider"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

var SetStop struct {
	SetStop struct {
		ID   string
		Name string
	} `graphql:"setStop(stop: $stop)"`
}

type Query struct {
	MapStopUpdated struct {
		ID       string
		Location Location
	} `graphql:"mapStopUpdated(minLat: $minLat, minLng: $minLng, maxLat: $maxLat, maxLng: $maxLng)"`
}

func TestRoundTrip(t *testing.T) {
	// m := sync.Mutex{}
	testSets := []struct {
		ID        string
		Latitude  float64
		Longitude float64
		StartTime time.Time
	}{
		{"123", 54.31981897337084, 10.182968719044112, time.Now()},
		{"last", 54.31981897337084, 10.182968719044112, time.Now()},
	}

	go func() {
		time.Sleep(1 * time.Second)

		client := graphql.NewClient("http://localhost:4567/query", nil).WithDebug(true)

		for _, testSet := range testSets {
			var m struct {
				SetStop struct {
					ID string
				} `graphql:"setStop(stop: $stop)"`
			}
			vars := map[string]interface{}{
				"stop": StopInput{
					ID:       testSet.ID,
					Provider: "test",
					Name:     "TestStop",
					Type:     "bus-stop",
					Location: Location{
						Longitude: testSet.Longitude,
						Latitude:  testSet.Latitude,
					},
				},
			}
			testSet.StartTime = time.Now()
			err := client.Mutate(context.Background(), &m, vars)
			assert.NoError(t, err)

			fmt.Println("mutate", testSet.ID, err)
		}
	}()

	client := graphql.NewSubscriptionClient("ws://localhost:4567/query").WithLog(log.Println)

	vars := map[string]interface{}{
		"minLat": 54.526130648172995,
		"minLng": 9.876994965672509,
		"maxLat": 53.95617973610979,
		"maxLng": 10.709999024470449,
	}
	q := Query{}
	_, err := client.Subscribe(&q, vars, func(dataValue []byte, err error) error {
		if err != nil {
			return err
		}

		data := Query{}
		// use the github.com/hasura/go-graphql-client/pkg/jsonutil package
		err = jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			return err
		}

		fmt.Println("updated", data.MapStopUpdated.ID)

		for _, testSet := range testSets {
			if data.MapStopUpdated.ID == testSet.ID {
				assert.Equal(t, testSet.Longitude, data.MapStopUpdated.Location.Longitude)
				assert.Equal(t, testSet.Latitude, data.MapStopUpdated.Location.Latitude)
				timeDiff := time.Since(testSet.StartTime)
				fmt.Println("round-trip", data.MapStopUpdated.ID, timeDiff)
				break
			}
		}

		if data.MapStopUpdated.ID == "last" {
			return graphql.ErrSubscriptionStopped
		}

		return nil
	})
	assert.NoError(t, err)
	defer client.Close()

	err = client.Run()
	assert.NoError(t, err)
}

package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
)

type graph struct{}

func (g *graph) Name() string {
	return "graphql"
}

func (g *graph) SendData(testSet *TestSet) error {
	client := graphql.NewClient("http://localhost:4567/query", nil)

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
	// fmt.Println("sending", testSet.ID)
	return client.Mutate(context.Background(), &m, vars)
}

func (g *graph) WaitForMessage(testSets []*TestSet, connectingWG *sync.WaitGroup, done func(s string)) error {
	client := graphql.NewSubscriptionClient("ws://localhost:4567/query")
	// WithLog(log.Println).
	// WithoutLogTypes(graphql.GQLData, graphql.GQLConnectionKeepAlive).
	// OnError(func(sc *graphql.SubscriptionClient, err error) error {
	// 	log.Print("err", err)
	// 	return err
	// })
	defer client.Close()

	vars := map[string]interface{}{
		"minLat": 54.526130648172995,
		"minLng": 9.876994965672509,
		"maxLat": 53.95617973610979,
		"maxLng": 10.709999024470449,
	}
	q := Query{}

	sub, err := client.Subscribe(&q, vars, func(dataValue []byte, err error) error {
		if err != nil {
			return err
		}

		data := Query{}
		// use the github.com/hasura/go-graphql-client/pkg/jsonutil package
		err = jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			return err
		}

		// only check round-trip for the first client
		for _, testSet := range testSets {
			if data.MapStopUpdated.ID == testSet.ID {
				if testSet.Longitude != data.MapStopUpdated.Location.Longitude || testSet.Latitude != data.MapStopUpdated.Location.Latitude {
					return fmt.Errorf("location mismatch: expected %f,%f, got %f,%f", testSet.Longitude, testSet.Latitude, data.MapStopUpdated.Location.Longitude, data.MapStopUpdated.Location.Latitude)
				}
				done(testSet.ID)
				break
			}
		}

		if data.MapStopUpdated.ID == "last" {
			return graphql.ErrSubscriptionStopped
		}

		return nil
	})
	if err != nil {
		return err
	}

	client.OnConnected(func() {
		// wait for the subscription to be running
		for {
			s := client.GetSubscription(sub)
			if s != nil && s.GetStatus() == graphql.SubscriptionRunning {
				// fmt.Printf("[%d:%d] connected\n", i, id)
				break
			}
		}

		connectingWG.Done()
	})

	err = client.Run()
	if err != nil && err.Error() != "exit" {
		return err
	}

	return nil
}

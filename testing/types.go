package testing

import (
	"sync"
	"time"
)

type Poc interface {
	Name() string
	SendData(testSet *TestSet) error
	WaitForMessage(testSets []*TestSet, connectingWG *sync.WaitGroup, done func(s string)) error
	SendPerf(amountClients string) error
}

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

type TestSet struct {
	ID        string
	Latitude  float64
	Longitude float64
	StartTime time.Time
}

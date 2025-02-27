package collector

import (
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
)

type Collector interface {
	Run(IDs []string)
	RunSingle(ID string)
	SubjectToID(subject string) string
	SubjectsToIDs(subjects []string) []string
}

func NewCollector(client *client.Client, collectorType string, subscriptions *subscriptions.Subscriptions) (Collector, error) {
	switch collectorType {
	case "map-vehicles":
		return &VehicleCollector{
			client: client,
		}, nil
	case "map-stops":
		return &StopCollector{
			client:        client,
			subscriptions: subscriptions,
		}, nil
	case "trips":
		return &TripCollector{
			client:        client,
			subscriptions: subscriptions,
		}, nil
	}

	return nil, fmt.Errorf("Collector type not supported")
}

package collector

import (
	"fmt"

	"github.com/kiel-live/kiel-live/pkg/client"
)

type Collector interface {
	Run()
	RunSingle(ID string)
	SubjectToID(subject string) string
}

func NewCollector(client client.Client, collectorType string) (Collector, error) {
	switch collectorType {
	case "vehicles":
		return &VehicleCollector{
			client: client,
		}, nil
	case "stops":
		return &StopCollector{
			client: client,
		}, nil
	case "trips":
		return &TripCollector{
			client: client,
		}, nil
	}

	return nil, fmt.Errorf("Collector type not supported")
}

package collector

import (
	"fmt"

	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

type Collector interface {
	Run()
	RunSingle(ID string)
	Reset()
	TopicToID(topic string) string
}

func NewCollector(client client.Client, collectorType string) (Collector, error) {
	switch collectorType {
	case "vehicles":
		return &VehicleCollector{
			client:   client,
			vehicles: make(map[string]*models.Vehicle),
		}, nil
	case "stops":
		return &StopCollector{
			client:         client,
			stops:          make(map[string]*models.Stop),
			lastFullUpdate: 0,
		}, nil
	case "trips":
		return &TripCollector{
			client: client,
			trips:  make(map[string]*models.Trip),
		}, nil
	}

	return nil, fmt.Errorf("Collector type not supported")
}

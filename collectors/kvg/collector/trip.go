package collector

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

type TripCollector struct {
	client client.Client
	trips  map[string]*models.Trip
	sync.Mutex
}

func isSameTripDepartures(a1, a2 []*models.TripDeparture) bool {
	if len(a1) != len(a2) {
		return false
	}
	if a1 == nil && a2 != nil || a1 != nil && a2 == nil {
		return false
	}
	for i, v := range a1 {
		if v != a2[i] {
			return false
		}
	}
	return true
}

func isSameTrip(a, b *models.Trip) bool {
	return a != nil && b != nil && a.ID == b.ID && a.Provider == b.Provider && a.Direction == b.Direction && isSameTripDepartures(a.Departures, b.Departures)
}

// returns list of changed or newly added trips
func (c *TripCollector) getChangedTrips(trips map[string]*models.Trip) (changed []*models.Trip) {
	for _, v := range trips {
		if !isSameTrip(v, c.trips[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *TripCollector) getRemovedTrips(trips map[string]*models.Trip) (removed []*models.Trip) {
	for _, v := range c.trips {
		if _, ok := trips[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *TripCollector) TopicToID(topic string) string {
	if strings.HasPrefix(topic, fmt.Sprintf(models.TopicTrip, api.IDPrefix)) && topic != fmt.Sprintf(models.TopicTrip, ">") {
		return strings.TrimPrefix(topic, fmt.Sprintf(models.TopicTrip, api.IDPrefix))
	}
	return ""
}

func (c *TripCollector) Run() {
	log := slog.With("collector", "trip")
	trips := map[string]*models.Trip{}

	c.Lock()
	defer c.Unlock()

	topics := c.client.GetSubscribedTopics()
	tripIDs := []string{}
	for _, topic := range topics {
		tripID := c.TopicToID(topic)
		if tripID != "" {
			tripIDs = append(tripIDs, tripID)
		}
	}

	for _, tripID := range tripIDs {
		trip, err := api.GetTrip(tripID)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		trips[trip.ID] = trip
	}

	// publish all changed trips
	changed := c.getChangedTrips(trips)
	for _, trip := range changed {
		log.Debug("publish changed trip", "trip", trip)
		err := c.client.UpdateTrip(trip)
		if err != nil {
			log.Error(err.Error())
		}
	}

	removed := c.getRemovedTrips(trips)
	for _, trip := range removed {
		log.Debug("publish removed trip", "trip", trip)
		err := c.client.DeleteTrip(trip.ID)
		if err != nil {
			log.Error(err.Error())
		}
	}

	log.Debug("collector run complete", "changed", len(changed), "removed", len(removed))
	// update list of trips
	c.trips = trips
}

func (c *TripCollector) RunSingle(tripID string) {
	log := slog.With("collector", "trip", "trip_id", tripID)

	c.Lock()
	defer c.Unlock()

	trip, err := api.GetTrip(tripID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// publish changed trip
	err = c.client.UpdateTrip(trip)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Debug("published single trip", "trip", trip)
	// update cache
	c.trips[trip.ID] = trip
}

func (c *TripCollector) Reset() {
	c.Lock()
	defer c.Unlock()

	c.trips = make(map[string]*models.Trip)
}

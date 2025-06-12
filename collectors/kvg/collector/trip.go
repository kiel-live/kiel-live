package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/sirupsen/logrus"
)

type TripCollector struct {
	client client.Client
	trips  map[string]*protocol.Trip
	sync.Mutex
}

func isSameTripArrivals(a1, a2 []protocol.TripArrival) bool {
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

func isSameTrip(a, b *protocol.Trip) bool {
	return a != nil && b != nil && a.ID == b.ID && a.Provider == b.Provider && a.Direction == b.Direction && isSameTripArrivals(a.Arrivals, b.Arrivals)
}

// returns list of changed or newly added trips
func (c *TripCollector) getChangedTrips(trips map[string]*protocol.Trip) (changed []*protocol.Trip) {
	for _, v := range trips {
		if !isSameTrip(v, c.trips[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *TripCollector) getRemovedTrips(trips map[string]*protocol.Trip) (removed []*protocol.Trip) {
	for _, v := range c.trips {
		if _, ok := trips[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *TripCollector) publish(trip *protocol.Trip) error {
	topic := fmt.Sprintf(protocol.TopicDetailsTrip, trip.ID)

	jsonData, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	err = c.client.Publish(topic, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (c *TripCollector) publishRemoved(trip *protocol.Trip) error {
	topic := fmt.Sprintf(protocol.TopicDetailsTrip, trip.ID)

	err := c.client.Publish(topic, string(protocol.DeletePayload))
	if err != nil {
		return err
	}

	return nil
}

func (c *TripCollector) TopicToID(topic string) string {
	if strings.HasPrefix(topic, fmt.Sprintf(protocol.TopicDetailsTrip, api.IDPrefix)) && topic != fmt.Sprintf(protocol.TopicDetailsTrip, ">") {
		return strings.TrimPrefix(topic, fmt.Sprintf(protocol.TopicDetailsTrip, api.IDPrefix))
	}
	return ""
}

func (c *TripCollector) Run() {
	log := logrus.WithField("collector", "trip")
	trips := map[string]*protocol.Trip{}

	c.Lock()
	defer c.Unlock()

	topics := c.subscriptions.GetSubscriptions()
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
			log.Error(err)
			continue
		}
		trips[trip.ID] = trip
	}

	// publish all changed trips
	changed := c.getChangedTrips(trips)
	for _, trip := range changed {
		log.Debugf("publish changed trip: %v", trip)
		err := c.publish(trip)
		if err != nil {
			log.Error(err)
		}
	}

	removed := c.getRemovedTrips(trips)
	for _, trip := range removed {
		log.Debugf("publish removed trip: %v", trip)
		err := c.publishRemoved(trip)
		if err != nil {
			log.Error(err)
		}
	}

	log.Debugf("changed %d trips and removed %d", len(changed), len(removed))
	// update list of trips
	c.trips = trips
}

func (c *TripCollector) RunSingle(tripID string) {
	log := logrus.WithField("collector", "trip").WithField("trip-id", tripID)

	c.Lock()
	defer c.Unlock()

	trip, err := api.GetTrip(tripID)
	if err != nil {
		log.Error(err)
		return
	}

	// publish changed trip
	err = c.publish(trip)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("published single trip: %v", trip)
	// update cache
	c.trips[trip.ID] = trip
}

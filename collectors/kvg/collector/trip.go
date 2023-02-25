package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/sirupsen/logrus"
)

type TripCollector struct {
	client        *client.Client
	trips         map[string]*protocol.Trip
	subscriptions *subscriptions.Subscriptions
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
	subject := fmt.Sprintf(protocol.SubjectDetailsTrip, trip.ID)

	jsonData, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	err = c.client.Publish(subject, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (c *TripCollector) publishRemoved(trip *protocol.Trip) error {
	subject := fmt.Sprintf(protocol.SubjectDetailsTrip, trip.ID)

	err := c.client.Publish(subject, string(protocol.DeletePayload))
	if err != nil {
		return err
	}

	return nil
}

func (c *TripCollector) SubjectsToIDs(subjects []string) []string {
	ids := []string{}
	for _, subject := range subjects {
		if strings.HasPrefix(subject, fmt.Sprintf(protocol.SubjectDetailsTrip, "")) && subject != fmt.Sprintf(protocol.SubjectDetailsTrip, ">") {
			ids = append(ids, strings.TrimPrefix(subject, fmt.Sprintf(protocol.SubjectDetailsTrip, "")+api.IDPrefix))
		}
	}
	return ids
}

func (c *TripCollector) Run(tripIDs []string, runRemoved bool) {
	log := logrus.WithField("collector", "trip")
	trips := map[string]*protocol.Trip{}
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

	var removed []*protocol.Trip
	if runRemoved {
		// publish all removed trips
		removed = c.getRemovedTrips(trips)
		for _, trip := range removed {
			log.Debugf("publish removed trip: %v", trip)
			err := c.publishRemoved(trip)
			if err != nil {
				log.Error(err)
			}
		}
	}

	log.Debugf("changed %d trips and removed %d", len(changed), len(removed))
	// update list of trips
	c.trips = trips
}

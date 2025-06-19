package collector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/sirupsen/logrus"
)

type StopCollector struct {
	client client.Client
	stops  map[string]*models.Stop
	sync.Mutex
}

func isSameArrivals(a, b []*models.StopArrival) bool {
	if len(a) != len(b) {
		return false
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func isSameStop(a *models.Stop, b *models.Stop) bool {
	return a != nil && b != nil &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.ID == b.ID &&
		isSameLocation(a.Location, b.Location) &&
		a.Type == b.Type &&
		isSameArrivals(a.Arrivals, b.Arrivals)
}

// returns list of changed or newly added stops
func (c *StopCollector) getChangedStops(stops map[string]*models.Stop) (changed []*models.Stop) {
	for _, v := range stops {
		if !isSameStop(v, c.stops[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *StopCollector) getRemovedStops(stops map[string]*models.Stop) (removed []*models.Stop) {
	for _, v := range c.stops {
		if _, ok := stops[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *StopCollector) TopicToID(topic string) string {
	if strings.HasPrefix(topic, fmt.Sprintf(models.TopicStop, api.IDPrefix)) && topic != fmt.Sprintf(models.TopicStop, ">") {
		return strings.TrimPrefix(topic, fmt.Sprintf(models.TopicStop, api.IDPrefix))
	}
	return ""
}

func (c *StopCollector) Run() {
	log := logrus.WithField("collector", "stop")

	c.Lock()
	defer c.Unlock()

	topics := c.client.GetSubscribedTopics()
	var stopIDs []string
	for _, topic := range topics {
		id := c.TopicToID(topic)
		if id != "" {
			stopIDs = append(stopIDs, id)
		}
	}

	stops, err := api.GetStops()
	if err != nil {
		log.Error(err)
		return
	}

	// load further details only for explicitly subscribed stops
	log.Debugf("loading details for %d stops: %s", len(stopIDs), strings.Join(stopIDs, ", "))
	for _, stopID := range stopIDs {
		details, err := api.GetStopDetails(stopID)
		if err != nil {
			log.Error(err)
			continue
		}
		stops[api.IDPrefix+stopID].Arrivals = details.Departures
		stops[api.IDPrefix+stopID].Alerts = details.Alerts
	}

	stopsToPublish := c.getChangedStops(stops)
	for _, stop := range stopsToPublish {
		log.Tracef("publish updated stop: %v", stop)
		err := c.client.UpdateStop(stop)
		if err != nil {
			log.Error(err)
		}
	}

	// publish all removed stops
	removed := c.getRemovedStops(stops)
	for _, stop := range removed {
		log.Debugf("publish removed stop: %v", stop)
		err := c.client.DeleteStop(stop.ID)
		if err != nil {
			log.Error(err)
		}
	}

	log.Debugf("changed %d stops and removed %d", len(stopsToPublish), len(removed))

	// update list of stops
	c.stops = stops
}

func (c *StopCollector) RunSingle(stopID string) {
	log := logrus.WithField("collector", "stop").WithField("stop-id", stopID)

	c.Lock()
	defer c.Unlock()

	// get stop from cache
	stop := c.stops[api.IDPrefix+stopID]

	// get stop details
	details, err := api.GetStopDetails(stopID)
	if err != nil {
		log.Error(err)
		return
	}
	stop.Arrivals = details.Departures
	stop.Alerts = details.Alerts

	// publish stop
	err = c.client.UpdateStop(stop)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("published single stop: %v", stop)
	// update cache
	c.stops[api.IDPrefix+stopID] = stop
}

func (c *StopCollector) Reset() {
	c.Lock()
	defer c.Unlock()

	c.stops = make(map[string]*models.Stop)
}

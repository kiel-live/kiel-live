package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/sirupsen/logrus"
)

type StopCollector struct {
	client         *client.Client
	stops          map[string]*protocol.Stop
	subscriptions  *subscriptions.Subscriptions
	lastFullUpdate int64
	sync.Mutex
}

func isSameArrivals(a, b []protocol.StopArrival) bool {
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

func isSameStop(a *protocol.Stop, b *protocol.Stop) bool {
	return a != nil && b != nil &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.ID == b.ID &&
		isSameLocation(a.Location, b.Location) &&
		a.Type == b.Type &&
		isSameArrivals(a.Arrivals, b.Arrivals)
}

// returns list of changed or newly added stops
func (c *StopCollector) getChangedStops(stops map[string]*protocol.Stop) (changed []*protocol.Stop) {
	for _, v := range stops {
		if !isSameStop(v, c.stops[v.ID]) {
			changed = append(changed, v)
		}
	}

	return changed
}

func (c *StopCollector) getRemovedStops(stops map[string]*protocol.Stop) (removed []*protocol.Stop) {
	for _, v := range c.stops {
		if _, ok := stops[v.ID]; !ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func (c *StopCollector) publish(stop *protocol.Stop) error {
	topic := fmt.Sprintf(protocol.TopicMapStop, stop.ID)

	jsonData, err := json.Marshal(stop)
	if err != nil {
		return err
	}

	err = c.client.Publish(topic, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (c *StopCollector) publishRemoved(stop *protocol.Stop) error {
	topic := fmt.Sprintf(protocol.TopicMapStop, stop.ID)

	err := c.client.Publish(topic, string(protocol.DeletePayload))
	if err != nil {
		return err
	}

	return nil
}

func (c *StopCollector) TopicToID(topic string) string {
	if strings.HasPrefix(topic, fmt.Sprintf(protocol.TopicMapStop, api.IDPrefix)) && topic != fmt.Sprintf(protocol.TopicMapStop, ">") {
		return strings.TrimPrefix(topic, fmt.Sprintf(protocol.TopicMapStop, api.IDPrefix))
	}
	return ""
}

func (c *StopCollector) Run() {
	log := logrus.WithField("collector", "stop")

	c.Lock()
	defer c.Unlock()

	topics := c.subscriptions.GetSubscriptions()
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

	for _, stopID := range stopIDs {
		log.Debug("StopCollector: Run: ", stopID)
		details, err := api.GetStopDetails(stopID)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Debug("StopCollector: publish stop", details)
		stops[api.IDPrefix+stopID].Arrivals = details.Departures
		stops[api.IDPrefix+stopID].Alerts = details.Alerts
	}

	var stopsToPublish []*protocol.Stop
	// publish all stops when last full update is older than the max cache age
	if c.lastFullUpdate == 0 || c.lastFullUpdate < time.Now().Unix()-protocol.MaxCacheAge {
		for _, stop := range stops {
			stopsToPublish = append(stopsToPublish, stop)
		}
		c.lastFullUpdate = time.Now().Unix()
	} else {
		// publish all changed stops
		stopsToPublish = c.getChangedStops(stops)
	}
	for _, stop := range stopsToPublish {
		err := c.publish(stop)
		if err != nil {
			log.Error(err)
		}
	}

	// publish all removed stops
	removed := c.getRemovedStops(stops)
	for _, stop := range removed {
		err := c.publishRemoved(stop)
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
	err = c.publish(stop)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("published single stop: %v", stop)
	// update cache
	c.stops[api.IDPrefix+stopID] = stop
}

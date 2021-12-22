package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

type StopCollector struct {
	client         *client.Client
	stops          map[string]*protocol.Stop
	subscriptions  *subscriptions.Subscriptions
	lastFullUpdate int64
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
	subject := fmt.Sprintf(protocol.SubjectMapStop, stop.ID)

	jsonData, err := json.Marshal(stop)
	if err != nil {
		return err
	}

	err = c.client.Publish(subject, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (c *StopCollector) publishRemoved(stop *protocol.Stop) error {
	subject := fmt.Sprintf(protocol.SubjectMapStop, stop.ID)

	err := c.client.Publish(subject, string(protocol.DeletePayload))
	if err != nil {
		return err
	}

	return nil
}

func (c *StopCollector) Run() {
	stops, err := api.GetStops()
	if err != nil {
		log.Error(err)
		return
	}

	for _, subject := range c.subscriptions.GetSubscriptions() {
		if !strings.HasPrefix(subject, fmt.Sprintf(protocol.SubjectMapStop, "")) || subject == fmt.Sprintf(protocol.SubjectMapStop, ">") {
			continue
		}
		// trim prefix of subject
		stopID := strings.TrimPrefix(subject, fmt.Sprintf(protocol.SubjectMapStop, "")+api.IDPrefix)
		log.Debug("StopCollector: Run: ", stopID)
		departures, err := api.GetStopDepartures(stopID)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Debug("StopCollector: publish stop", departures)
		stops[api.IDPrefix+stopID].Arrivals = departures
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
		c.publish(stop)
	}

	// publish all removed stops
	removed := c.getRemovedStops(stops)
	for _, stop := range removed {
		c.publishRemoved(stop)
	}

	log.Debugf("changed %d stops and removed %d", len(stopsToPublish), len(removed))

	// update list of stops
	c.stops = stops
}

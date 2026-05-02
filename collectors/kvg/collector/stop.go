package collector

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

type StopCollector struct {
	client client.Client
	stops  map[string]*models.Stop
	sync.Mutex
	lastFullUpdate int64
}

func isSameDepartures(a, b []*models.StopDepartures) bool {
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
		isSameDepartures(a.Departures, b.Departures)
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
	log := slog.With("collector", "stop")

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
		log.Error("failed to get stops", "error", err)
		return
	}

	// load further details only for explicitly subscribed stops
	log.Debug("loading details for stops", "count", len(stopIDs), "stop_ids", strings.Join(stopIDs, ", "))
	for _, stopID := range stopIDs {
		details, err := api.GetStopDetails(stopID)
		if err != nil {
			log.Error("failed to get stop details", "error", err)
			continue
		}
		stops[api.IDPrefix+stopID].Departures = details.Departures
		stops[api.IDPrefix+stopID].Alerts = details.Alerts
	}

	var stopsToPublish []*models.Stop
	// publish all stops when last full update is older than the max cache age as stops
	// normally never change and the cache wont keep them forever
	if c.lastFullUpdate == 0 || c.lastFullUpdate < time.Now().Unix()-models.MaxCacheAge {
		for _, stop := range stops {
			stopsToPublish = append(stopsToPublish, stop)
		}
		c.lastFullUpdate = time.Now().Unix()
	} else {
		// publish all changed stops
		stopsToPublish = c.getChangedStops(stops)
	}
	for _, stop := range stopsToPublish {
		log.Debug("publish updated stop", "stop", stop)
		err := c.client.UpdateStop(stop)
		if err != nil {
			log.Error("failed to update stop", "error", err)
		}
	}

	// publish all removed stops
	removed := c.getRemovedStops(stops)
	for _, stop := range removed {
		log.Debug("publish removed stop", "stop", stop)
		err := c.client.DeleteStop(stop.ID)
		if err != nil {
			log.Error("failed to delete stop", "error", err)
		}
	}

	log.Debug("collector run complete", "changed", len(stopsToPublish), "removed", len(removed))

	// update list of stops
	c.stops = stops
}

func (c *StopCollector) RunSingle(stopID string) {
	log := slog.With("collector", "stop", "stop_id", stopID)

	c.Lock()
	defer c.Unlock()

	// get stop from cache
	stop, ok := c.stops[api.IDPrefix+stopID]
	if !ok {
		log.Debug("stop not found in cache, fetching stops list")
		stops, err := api.GetStops()
		if err != nil {
			log.Error("failed to get stops", "error", err)
			return
		}

		stop, ok = stops[api.IDPrefix+stopID]
		if !ok {
			log.Error("stop not found in stops list")
			return
		}
	}

	// get stop details
	details, err := api.GetStopDetails(stopID)
	if err != nil {
		log.Error("failed to get stop details", "error", err)
		return
	}
	stop.Departures = details.Departures
	stop.Alerts = details.Alerts

	// publish stop
	err = c.client.UpdateStop(stop)
	if err != nil {
		log.Error("failed to update stop", "error", err)
	}

	log.Debug("published single stop", "stop", stop)
	// update cache
	c.stops[api.IDPrefix+stopID] = stop
}

func (c *StopCollector) Reset() {
	c.Lock()
	defer c.Unlock()

	c.stops = make(map[string]*models.Stop)
}

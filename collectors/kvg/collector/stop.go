package collector

import (
	"encoding/json"
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/api"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

type StopCollector struct {
	client *client.Client
	stops  map[string]*protocol.Stop
}

func isSameStop(a *protocol.Stop, b *protocol.Stop) bool {
	return a != nil && b != nil &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.ID == b.ID &&
		isSameLocation(a.Location, b.Location) &&
		a.Type == b.Type
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

	err := c.client.Publish(subject, string("---"))
	if err != nil {
		return err
	}

	return nil
}

func (c *StopCollector) Run() {
	stops := api.GetStops()

	// publish all changed stops
	changed := c.getChangedStops(stops)
	for _, stop := range changed {
		c.publish(stop)
	}

	// publish all removed stops
	removed := c.getRemovedStops(stops)
	for _, stop := range removed {
		c.publishRemoved(stop)
	}

	log.Debugf("changed %d stops and removed %d", len(changed), len(removed))

	// update list of stops
	c.stops = stops
}

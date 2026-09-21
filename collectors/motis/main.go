package main

import (
	"os"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/collectors/motis/api"
	"github.com/kiel-live/kiel-live/collectors/motis/version"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
	log "github.com/sirupsen/logrus"
)

// Collector manages MOTIS stop data. It periodically refreshes the full stop
// list and, for stops that clients have subscribed to, also fetches departures
// on demand and on every subsequent tick.
type Collector struct {
	client          client.Client
	stops           map[string]*api.StopWithProviderID // all known stops from MOTIS
	subscribedStops map[string]bool                    // stop IDs currently subscribed by clients
	lastStopsFetch  int64
	sync.Mutex
}

// refreshStops fetches the full stop list from MOTIS, publishes all stops,
// and removes stops that have disappeared from the feed.
// Must be called with the lock held.
func (c *Collector) refreshStops() {
	log := log.WithField("collector", "motis")

	stops, err := api.GetStops()
	if err != nil {
		log.Errorf("failed to fetch stops: %v", err)
		return
	}

	for id, stop := range c.stops {
		if _, ok := stops[id]; !ok {
			if err := c.client.DeleteStop(stop.ID); err != nil {
				log.Errorf("failed to delete stop %s: %v", stop.ID, err)
			}
		}
	}

	for _, stop := range stops {
		if err := c.client.UpdateStop(stop.Stop); err != nil {
			log.Errorf("failed to publish stop %s: %v", stop.ID, err)
		}
	}

	c.stops = stops
	c.lastStopsFetch = time.Now().Unix()
	log.Debugf("fetched and published %d stops", len(stops))
}

// refreshDepartures fetches departures for a single stop and publishes the
// updated stop. Must be called with the lock held.
func (c *Collector) refreshDepartures(stopID string) {
	log := log.WithField("collector", "motis").WithField("stop-id", stopID)

	stop, ok := c.stops[stopID]
	if !ok {
		log.Debugf("stop not in cache, skipping departures")
		return
	}

	departures, err := api.GetStopDepartures(stop.ProviderIDs)
	if err != nil {
		log.Errorf("failed to fetch departures: %v", err)
		return
	}

	// Clone stop and attach departures so we don't mutate the cached entry.
	updated := *stop
	updated.Departures = departures
	// Explicit empty slice so consumers can distinguish "no departures" from "not loaded yet".
	if len(departures) == 0 {
		updated.Departures = make([]*models.StopDepartures, 0)
	}

	if err := c.client.UpdateStop(updated.Stop); err != nil {
		log.Errorf("failed to publish stop with departures: %v", err)
		return
	}

	log.Debugf("published %d departures", len(departures))
}

// tick is called on every scheduler interval. It keeps the stop list fresh
// and refreshes departures for all currently subscribed stops.
func (c *Collector) tick() {
	c.Lock()
	defer c.Unlock()

	if c.lastStopsFetch == 0 || time.Now().Unix()-c.lastStopsFetch > models.MaxCacheAge {
		c.refreshStops()
	}

	for stopID := range c.subscribedStops {
		c.refreshDepartures(stopID)
	}
}

// subscribe adds a stop to the subscribed set and immediately fetches its
// departures so the client does not have to wait for the next tick.
func (c *Collector) subscribe(stopID string) {
	c.Lock()
	defer c.Unlock()

	if c.stops == nil {
		c.refreshStops()
	}

	c.subscribedStops[stopID] = true
	c.refreshDepartures(stopID)
}

// unsubscribe removes a stop from the subscribed set. Departures for this
// stop will no longer be refreshed on subsequent ticks.
func (c *Collector) unsubscribe(stopID string) {
	c.Lock()
	defer c.Unlock()

	delete(c.subscribedStops, stopID)
}

// reset clears the stop list cache, forcing a full refresh on the next tick.
// Subscribed stops are retained so departure refreshes resume automatically
// after reconnect.
func (c *Collector) reset() {
	c.Lock()
	defer c.Unlock()

	c.stops = nil
	c.lastStopsFetch = 0
}

func main() {
	log.Infof("Kiel-Live MOTIS collector version %s", version.Version)

	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Fatalf("error loading location '%s': %v\n", tz, err)
		}
	}

	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	server := os.Getenv("COLLECTOR_SERVER")
	if server == "" {
		log.Fatalln("Please provide a server address for the collector with COLLECTOR_SERVER")
	}

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the collector with COLLECTOR_TOKEN")
	}

	c := client.NewClient(server, token)
	err = c.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer func() {
		if err := c.Disconnect(); err != nil {
			log.Error(err)
		}
	}()

	collector := &Collector{
		client:          c,
		subscribedStops: make(map[string]bool),
	}

	c.SetOnTopicsChanged(func(topic string, subscribed bool) {
		id := api.TopicToID(topic)
		if id == "" {
			return
		}

		if subscribed {
			collector.subscribe(id)
		} else {
			collector.unsubscribe(id)
		}
	})

	c.SetOnConnectionChanged(func(connected bool) {
		if !connected {
			return
		}
		log.Debug("Reconnected, resetting stop cache")
		collector.reset()
	})

	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Error(err)
		}
	}()

	_, err = s.NewJob(
		gocron.DurationJob(30*time.Second),
		gocron.NewTask(func() {
			if !c.IsConnected() {
				return
			}
			collector.tick()
		}),
	)
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("MOTIS collector started")
	s.Start()
	select {}
}

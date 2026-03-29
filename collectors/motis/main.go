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

// StopCollector fetches non-bus stops and their departures from MOTIS.
//
// Stop metadata (name, location, type) is refreshed every MaxCacheAge seconds.
// Departures are fetched on demand when a stop is subscribed, and refreshed
// on every Run() for all currently subscribed stops.
type StopCollector struct {
	client         client.Client
	stops          map[string]*api.StopWithProviderID
	lastStopsFetch int64
	sync.Mutex
}

// refreshStops fetches the stop list from MOTIS and publishes all stops.
// Must be called with the lock held.
func (c *StopCollector) refreshStops() {
	log := log.WithField("collector", "stop")

	stops, err := api.GetStops()
	if err != nil {
		log.Errorf("failed to fetch stops: %v", err)
		return
	}

	// Detect and publish removed stops
	for id, stop := range c.stops {
		if _, ok := stops[id]; !ok {
			if err := c.client.DeleteStop(stop.ID); err != nil {
				log.Errorf("failed to delete stop %s: %v", stop.ID, err)
			}
		}
	}

	// Publish all stops (metadata only, no arrivals yet)
	for _, stop := range stops {
		if err := c.client.UpdateStop(stop.Stop); err != nil {
			log.Errorf("failed to publish stop %s: %v", stop.ID, err)
		}
	}

	c.stops = stops
	c.lastStopsFetch = time.Now().Unix()
	log.Debugf("fetched and published %d stops", len(stops))
}

// fetchDepartures loads departure data for a single stop and publishes it.
// Must be called with the lock held.
func (c *StopCollector) fetchDepartures(stopID string) {
	log := log.WithField("collector", "stop").WithField("stop-id", stopID)

	stop, ok := c.stops[stopID]
	if !ok {
		log.Debugf("stop not in cache, skipping departures")
		return
	}

	arrivals, err := api.GetStopDepartures(stop.ProviderIDs)
	if err != nil {
		log.Errorf("failed to fetch departures: %v", err)
		return
	}

	// Clone stop and attach arrivals before publishing so we don't mutate
	// the cached stop (arrivals change every run; metadata stays stable).
	updated := *stop
	updated.Arrivals = arrivals

	if err := c.client.UpdateStop(updated.Stop); err != nil {
		log.Errorf("failed to publish stop with departures: %v", err)
		return
	}

	log.Debugf("published %d departures", len(arrivals))
}

// Run is called on every scheduler tick and refreshes the stop list if the
// cache has expired. Departures are not fetched here — they are loaded on
// demand when a client subscribes to a stop topic (see RunSingle).
func (c *StopCollector) Run() {
	c.Lock()
	defer c.Unlock()

	if c.lastStopsFetch == 0 || time.Now().Unix()-c.lastStopsFetch > models.MaxCacheAge {
		c.refreshStops()
	}
}

// RunSingle is called immediately when a client subscribes to a stop topic.
// It ensures the stop (with fresh departures) is published without waiting
// for the next Run() tick.
func (c *StopCollector) RunSingle(stopID string) {
	c.Lock()
	defer c.Unlock()

	// Ensure the stop list is populated
	if c.stops == nil {
		c.refreshStops()
	}

	c.fetchDepartures(stopID)
}

// Reset clears the stop cache, forcing a full refresh on the next Run().
func (c *StopCollector) Reset() {
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

	stops := &StopCollector{client: c}

	c.SetOnTopicsChanged(func(topic string, subscribed bool) {
		if !subscribed {
			return
		}
		id := api.TopicToID(topic)
		if id == "" {
			return
		}

		stops.RunSingle(id)
	})

	c.SetOnConnectionChanged(func(connected bool) {
		if !connected {
			return
		}
		log.Debug("Reconnected, resetting stop collector")
		stops.Reset()
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
			stops.Run()
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

package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/collectors/kvg/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
	log "github.com/sirupsen/logrus"
)

var collectors map[string]collector.Collector

func main() {
	log.Infof("🚌 Kiel-Live KVG collector version %s", "1.0.0") // TODO use proper version

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
		err := c.Disconnect()
		if err != nil {
			log.Error(err)
		}
	}()

	collectors = make(map[string]collector.Collector)

	// auto load following collectors
	collectors["vehicles"], err = collector.NewCollector(c, "vehicles")
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["stops"], err = collector.NewCollector(c, "stops")
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["trips"], err = collector.NewCollector(c, "trips")
	if err != nil {
		log.Errorln(err)
		return
	}

	c.SetOnTopicsChanged(func(topic string, subscribed bool) {
		if !subscribed {
			return
		}

		tripID := collectors["trips"].TopicToID(topic)
		if tripID != "" {
			collectors["trips"].RunSingle(tripID)
			return
		}
		stopID := collectors["stops"].TopicToID(topic)
		if stopID != "" {
			collectors["stops"].RunSingle(stopID)
			return
		}
	})

	c.SetOnConnectionChanged(func(connected bool) {
		if !connected {
			return
		}

		for name, c := range collectors {
			log.Debugf("Resetting %s collector", name)
			c.Reset()
		}
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
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(func() {
			if !c.IsConnected() {
				return
			}

			for name, c := range collectors {
				// TODO maybe run in go routine
				log.Debugf("Running %s collector ...", name)
				c.Run()
			}
		}),
	)
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ KVG collector started")

	s.Start()
	select {}
}

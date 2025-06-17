package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
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

	c := client.NewNatsClient(server, client.WithAuth("collector", token))
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

	c.SetOnTopicsChanged(func(topic string, added bool) {
		if !added {
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

		// TODO: reset collector caches and re-run them
	})

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err = s.Every(5).Seconds().Do(func() {
		if !c.IsConnected() {
			return
		}

		for name, c := range collectors {
			// TODO maybe run in go routine
			log.Debugln("Collector for", name, "running ...")
			c.Run()
		}
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ KVG collector started")

	s.StartBlocking()
}

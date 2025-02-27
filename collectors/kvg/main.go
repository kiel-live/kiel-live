package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/collector"
	"github.com/kiel-live/kiel-live/collectors/kvg/subscriptions"
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

	c := client.NewClient(server, client.WithAuth("collector", token))
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

	subscriptions := subscriptions.New(c)

	collectors = make(map[string]collector.Collector)

	// auto load following collectors
	collectors["map-vehicles"], err = collector.NewCollector(c, "map-vehicles", subscriptions)
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["map-stops"], err = collector.NewCollector(c, "map-stops", subscriptions)
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["trips"], err = collector.NewCollector(c, "trips", subscriptions)
	if err != nil {
		log.Errorln(err)
		return
	}

	subscriptions.Subscribe(func(subject string) {
		tripID := collectors["trips"].SubjectToID(subject)
		if tripID != "" {
			collectors["trips"].RunSingle(tripID)
			return
		}
		stopID := collectors["map-stops"].SubjectToID(subject)
		if stopID != "" {
			collectors["map-stops"].RunSingle(stopID)
			return
		}
	})

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err = s.Every(5).Seconds().Do(func() {
		if !c.IsConnected() {
			return
		}

		subjects := subscriptions.GetSubscriptions()
		for name, c := range collectors {
			// TODO maybe run in go routine
			log.Debugln("Collector for", name, "running ...")
			c.Run(c.SubjectsToIDs(subjects))
		}
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ KVG collector started")

	s.StartBlocking()
}

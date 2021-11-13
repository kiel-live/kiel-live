package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/kvg/collector"
	"github.com/kiel-live/kiel-live/protocol"
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
		log.Fatalln("Please provide a token for the collector with MANAGER_TOKEN")
	}

	c := client.NewClient(server, client.WithAuth("collector", token))
	err = c.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer c.Disconnect()

	collectors = make(map[string]collector.Collector)

	// auto load following collectors
	collectors["map-vehicles"], err = collector.NewCollector(c, "map-vehicles")
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["map-stops"], err = collector.NewCollector(c, "map-stops")
	if err != nil {
		log.Errorln(err)
		return
	}

	err = c.Subscribe(protocol.SubjectSubscriptions, func(msg *client.SubjectMessage) {
		log.Debugln("Subscriptions", msg.Data)
		// TODO start & stop on-demand collectors
	}, c.WithCache())
	if err != nil {
		log.Errorln("Ups", err)
		return
	}

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.Every(1).Seconds().Do(func() {
		if !c.IsConnected() {
			return
		}

		for name, c := range collectors {
			// TODO maybe run in go routine
			log.Debugln("Collector for", name, "running ...")
			c.Run()
		}
	})

	log.Infoln("⚡ KVG collector started")

	s.StartBlocking()
}

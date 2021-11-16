package main

import (
	"encoding/json"
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

// keep list of consumers for easy deletion
var consumers map[string]string

// keep track of duplicate subscriptions
var subjects map[string]int

// plain list of subscriptions without duplicates
var subscriptions []string

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
	collectors["map-vehicles"], err = collector.NewCollector(c, "map-vehicles", &subscriptions)
	if err != nil {
		log.Errorln(err)
		return
	}
	collectors["map-stops"], err = collector.NewCollector(c, "map-stops", &subscriptions)
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

	// TODO: move to Subscriptions Handler package
	type consumerEvent struct {
		Stream   string `json:"stream"`
		Consumer string `json:"consumer"`
	}

	consumers = make(map[string]string)
	subjects = make(map[string]int)
	subscriptions = []string{}

	// already existing consumers
	for consumerInfo := range c.JS.ConsumersInfo("data") {
		consumers[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		subjects[consumerInfo.Config.FilterSubject]++
		subscriptions = []string{}
		for subject := range subjects {
			subscriptions = append(subscriptions, subject)
		}
	}

	// new consumers
	c.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.CREATED.>", func(msg *client.SubjectMessage) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		consumerInfo, _ := c.JS.ConsumerInfo(consumerEvent.Stream, consumerEvent.Consumer)
		consumers[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		subjects[consumerInfo.Config.FilterSubject]++
		subscriptions = []string{}
		for subject := range subjects {
			subscriptions = append(subscriptions, subject)
		}
		log.Debugln("Subscriptions", consumers)
		log.Debugln("Subscriptions", subjects)
		log.Debugln("Subscriptions", subscriptions)
		collectors["map-stops"].Run()
	})

	// remove consumers
	c.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.DELETED.>", func(msg *client.SubjectMessage) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		if subjects[consumers[consumerEvent.Consumer]] > 1 {
			subjects[consumers[consumerEvent.Consumer]]--
		} else {
			delete(subjects, consumers[consumerEvent.Consumer])
		}
		delete(consumers, consumerEvent.Consumer)
		subscriptions = []string{}
		for subject := range subjects {
			subscriptions = append(subscriptions, subject)
		}
		log.Debugln("Subscriptions", consumers)
		log.Debugln("Subscriptions", subjects)
		log.Debugln("Subscriptions", subscriptions)
	})

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.Every(5).Seconds().Do(func() {
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

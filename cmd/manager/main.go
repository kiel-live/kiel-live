package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/manager"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("🚌 Kiel-Live manager version %s", "2.0.0") // TODO use proper version

	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	server := os.Getenv("MANAGER_SERVER")
	if server == "" {
		log.Fatalln("Please provide a server address for the manager with MANAGER_SERVER")
	}

	token := os.Getenv("MANAGER_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the manager with MANAGER_TOKEN")
	}

	c := client.NewClient(server, client.WithAuth("manager", token))
	err = c.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer c.Disconnect()

	hub := manager.NewHub()
	defer hub.Unload()

	err = c.Subscribe(protocol.SubjectRequestSubscribe, func(msg *client.SubjectMessage) {
		subject := string(msg.Data)
		log.Debugln("Subscribe to", subject)
		err := hub.Subscribe(subject)
		if err != nil {
			log.Errorln(err)
			msg.Raw.Respond([]byte("err"))
			return
		}

		err = msg.Raw.Respond([]byte("ok"))
		if err != nil {
			log.Errorln(err)
			return
		}
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = c.Subscribe(protocol.SubjectRequestUnsubscribe, func(msg *client.SubjectMessage) {
		subject := string(msg.Data)
		log.Debugln("Unsubscribe from", subject)
		err := hub.Unsubscribe(subject)

		if err != nil {
			log.Errorln(err)
			msg.Raw.Respond([]byte("err"))
			return
		}

		msg.Raw.Respond([]byte("ok"))
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	c.JS.AddStream(&nats.StreamConfig{
		Name:              "data",
		Subjects:          []string{"data.>"},
		Retention:         nats.LimitsPolicy,
		MaxMsgsPerSubject: 1,
		Discard:           nats.DiscardOld,
		MaxAge:            time.Duration(protocol.MaxCacheAge * time.Second),
	})

	log.Infoln("⚡ Manager started")

	// don't kill main process
	select {}
}

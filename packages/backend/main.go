package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/packages/backend/hub"
	"github.com/kiel-live/kiel-live/packages/backend/store"
	"github.com/kiel-live/kiel-live/packages/backend/webserver"
	"github.com/kiel-live/kiel-live/packages/backend/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "Main")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	port := 4000
	if _port, ok := os.LookupEnv("PORT"); ok {
		port, err = strconv.Atoi(_port)
		if err != nil {
			log.Panic("Please provide a port as number with PORT")
		}
	}

	collectorToken := os.Getenv("COLLECTOR_TOKEN")
	if collectorToken == "" {
		log.Panic("Please provide a token for collector access with COLLECTOR_TOKEN")
	}

	log.Printf("🚌 Kiel-Live backend version %s", "2.0.0") // TODO load proper version
	log.Println("⚡ Backend starting ...")

	store := store.NewMemoryStore()
	store.Load()
	defer store.Unload()

	hub, err := hub.NewHub(store)
	if err != nil {
		log.Panic("Can't start hub")
	}
	go hub.Run()

	websocketServer := websocket.NewServer(hub, collectorToken)

	webServer := webserver.NewWebServer(websocketServer)
	err = webServer.Listen(port)
	if err != nil {
		log.Panic("Can't start web-server")
	}

	defer webServer.Close()
}

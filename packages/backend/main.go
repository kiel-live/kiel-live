package main

import (
	"flag"
	"log"

	"github.com/kiel-live/kiel-live/packages/backend/hub"
	"github.com/kiel-live/kiel-live/packages/backend/store"
	"github.com/kiel-live/kiel-live/packages/backend/webserver"
	"github.com/kiel-live/kiel-live/packages/backend/websocket"
)

var port = flag.Int("addr", 4000, "server port")

func main() {
	flag.Parse()

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

	websocketServer := websocket.NewServer(hub)

	webServer := webserver.NewWebServer(websocketServer)
	err = webServer.Listen(*port)
	if err != nil {
		log.Panic("Can't start web-server")
	}

	defer webServer.Close()
}

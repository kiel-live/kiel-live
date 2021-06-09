package main

import (
	"flag"
	"log"

	"github.com/kiel-live/kiel-live/backend/hub"
	"github.com/kiel-live/kiel-live/backend/webserver"
	"github.com/kiel-live/kiel-live/backend/websocket"
)

var port = flag.Int("addr", 4000, "server port")

func main() {
	flag.Parse()

	log.Printf("🚌 Kiel-Live backend version %s", "2.0.0") // TODO load proper version
	log.Println("⚡ Backend starting ...")

	hub := hub.NewHub()
	go hub.Run()

	websocketServer := websocket.NewWebsocketServer(hub)

	webServer := webserver.NewWebServer(websocketServer)
	err := webServer.Listen(*port)
	if err != nil {
		log.Panic("Can't start web-server")
	}

	defer webServer.Close()
}

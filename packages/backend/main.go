package main

import (
	"flag"
	"log"

	"github.com/kiel-live/kiel-live/backend/hub"
	"github.com/kiel-live/kiel-live/backend/webserver"
	"github.com/kiel-live/kiel-live/backend/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	log.Println("Starting server ...")

	hub := hub.NewHub()
	go hub.Run()

	websocketServer := websocket.NewWebsocketServer(hub)

	webServer := webserver.NewWebServer(websocketServer)
	err := webServer.Listen(*addr)
	if err != nil {
		log.Panic("Can't start web-server")
	}

	defer webServer.Close()
}

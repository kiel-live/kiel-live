package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/rpc"
	"github.com/kiel-live/kiel-live/shared/database"
	"github.com/kiel-live/kiel-live/shared/hub"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

const defaultPort = "4567"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(_rpc *rpc.RPC, pub pubsub.Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = rpc.NewServer(_rpc, pub, conn.UnderlyingConn())
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.NewMemoryDatabase()
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	hub := &hub.Hub{
		DB:     db,
		PubSub: pubsub.NewMemory(),
	}
	rp := rpc.NewRPC(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(rp, hub.PubSub, w, r)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/kiel-live/kiel-live/jsonrpc/rpc"
	"github.com/kiel-live/kiel-live/shared/database"
	"github.com/kiel-live/kiel-live/shared/hub"
	"github.com/kiel-live/kiel-live/shared/pubsub"
)

const defaultPort = "4568"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(ctx context.Context, pub pubsub.Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	peer := rpc.NewServerPeer(ctx, conn.UnderlyingConn(), pub)

	err = peer.Register(&KielLiveRPC{})
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

	ctx := context.Background()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(ctx, hub.PubSub, w, r)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

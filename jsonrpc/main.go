package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"

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

	broker := pubsub.NewMemory()
	hub := &hub.Hub{
		DB:     db,
		PubSub: broker,
	}

	ctx := context.Background()

	server := rpc.NewServer(broker)
	err := server.Register(&KielLiveRPC{
		Hub: hub,
	})
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		server.NewPeer(ctx, websocketjsonrpc2.NewObjectStream(conn))
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

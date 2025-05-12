package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"

	"github.com/kiel-live/kiel-live/pkg/database"
	"github.com/kiel-live/kiel-live/pkg/pubsub"
	"github.com/kiel-live/kiel-live/pkg/rpc"
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
		http.Error(w, "not found", http.StatusNotFound)
	})

	broker := pubsub.NewMemory()
	server := rpc.NewServer(broker)
	err := server.RegisterName("hub", &Hub{
		DB:     db,
		PubSub: broker,
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

		ctx := context.Background()
		server.NewPeer(ctx, websocketjsonrpc2.NewObjectStream(conn))
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

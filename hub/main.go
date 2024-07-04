package main

import (
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/olahol/melody"

	"github.com/kiel-live/kiel-live/shared/database"
)

const defaultPort = "4567"

type GopherInfo struct {
	ID string
	X  string
	Y  string
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

	m := melody.New()

	// _rpc := rpc.RPC{
	// 	// hub: &hub.Hub{
	// 	// 	DB: db,
	// 	// },
	// }

	m.HandleConnect(func(s *melody.Session) {
		ss, _ := m.Sessions()

		for _, o := range ss {
			value, exists := o.Get("info")

			if !exists {
				continue
			}

			info := value.(*GopherInfo)

			s.Write([]byte("set " + info.ID + " " + info.X + " " + info.Y))
		}

		id := uuid.NewString()
		s.Set("info", &GopherInfo{id, "0", "0"})

		s.Write([]byte("iam " + id))

		// server, err := rpc.NewServer(_rpc, s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		value, exists := s.Get("info")

		if !exists {
			return
		}

		info := value.(*GopherInfo)

		m.BroadcastOthers([]byte("dis "+info.ID), s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		err := m.Broadcast(msg)
		if err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		err := m.HandleRequest(w, r)
		if err != nil {
			log.Println(err)
		}
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool { // Allow all connections
		return true
	},
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	s.hub.Register <- conn
	log.Println("WebSocket client connected")

	// Handle incoming messages from this client for subscriptions
	go func(c *websocket.Conn) {
		defer func() {
			s.hub.Unregister <- c
			c.Close()
			log.Println("WebSocket client disconnected and unregistered")
		}()
		for {
			var msg hub.ClientSubscriptionMessage
			if err := c.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error reading json from client: %v", err)
				}
				break // Exit loop on error
			}

			switch strings.ToLower(msg.Action) {
			case "subscribe":
				s.hub.Subscribe <- hub.SubscriptionRequest{Client: c, Topic: msg.Topic}
			case "unsubscribe":
				s.hub.Unsubscribe <- hub.SubscriptionRequest{Client: c, Topic: msg.Topic}
			default:
				log.Printf("Unknown action from client: %s", msg.Action)
			}
		}
	}(conn)
}

package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool { // Allow all connections
		return true
	},
}

// clientSubscriptionMessage defines the structure for client messages to subscribe/unsubscribe.
type clientSubscriptionMessage struct {
	Action string `json:"action"` // "subscribe", "unsubscribe", or "search"
	Topic  string `json:"topic"`
	Query  string `json:"query,omitempty"` // for search
	Limit  int    `json:"limit,omitempty"` // for search
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	s.hub.RegisterClient(conn)
	log.Println("WebSocket client connected")

	// Handle incoming messages from this client for subscriptions
	go func(c *websocket.Conn) {
		defer func() {
			s.hub.UnregisterClient(c)
			c.Close()
			log.Println("WebSocket client disconnected and unregistered")
		}()
		for {
			var msg clientSubscriptionMessage
			if err := c.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error reading json from client: %v", err)
				}
				break // Exit loop on error
			}

			switch strings.ToLower(msg.Action) {
			case "subscribe":
				s.hub.SubscribeClient(c, msg.Topic)
			case "unsubscribe":
				s.hub.UnsubscribeClient(c, msg.Topic)
			case "search":
				limit := msg.Limit
				if limit == 0 {
					limit = 20
				}
				results, err := s.search.Search(r.Context(), msg.Query, limit)
				if err != nil {
					log.Printf("search error: %v", err)
				} else {
					if err := c.WriteJSON(map[string]any{
						"action": "search_results",
						"data":   results,
					}); err != nil {
						log.Printf("error writing search results: %v", err)
					}
				}
			default:
				log.Printf("Unknown action from client: %s", msg.Action)
			}
		}
	}(conn)
}

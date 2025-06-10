package hub

import (
	"log"
	"sync"

	"github.com/kiel-live/kiel-live/pkg/database"

	"github.com/gorilla/websocket"
)

// WebsocketMessage defines the structure for messages sent over WebSocket.
type WebsocketMessage struct {
	Topic  string `json:"topic"`
	Action string `json:"action"` // e.g., "created", "updated", "deleted"
	Data   any    `json:"data"`
}

// ClientSubscriptionMessage defines the structure for client messages to subscribe/unsubscribe.
type ClientSubscriptionMessage struct {
	Action string `json:"action"` // "subscribe" or "unsubscribe"
	Topic  string `json:"topic"`
}

type SubscriptionRequest struct {
	Client *websocket.Conn
	Topic  string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients     map[*websocket.Conn]map[string]struct{} // client -> set of subscribed topics
	mu          sync.Mutex                              // To protect clients map
	Broadcast   chan WebsocketMessage
	Register    chan *websocket.Conn
	Unregister  chan *websocket.Conn
	Subscribe   chan SubscriptionRequest
	Unsubscribe chan SubscriptionRequest
	db          database.Database
}

func NewHub(db database.Database) *Hub {
	return &Hub{
		Broadcast:   make(chan WebsocketMessage),
		Register:    make(chan *websocket.Conn),
		Unregister:  make(chan *websocket.Conn),
		Subscribe:   make(chan SubscriptionRequest),
		Unsubscribe: make(chan SubscriptionRequest),
		clients:     make(map[*websocket.Conn]map[string]struct{}),
		db:          db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = make(map[string]struct{}) // Initialize empty set of topics
			h.mu.Unlock()
			log.Println("Client registered")
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
		case req := <-h.Subscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.Client]; ok {
				clientTopics[req.Topic] = struct{}{}
				log.Printf("Client subscribed to topic: %s", req.Topic)
			} else {
				log.Printf("Failed to subscribe: client not registered: %v", req.Client)
			}
			h.mu.Unlock()
		case req := <-h.Unsubscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.Client]; ok {
				delete(clientTopics, req.Topic)
				log.Printf("Client unsubscribed from topic: %s", req.Topic)
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			for client, topics := range h.clients {
				if _, subscribed := topics[message.Topic]; subscribed {
					err := client.WriteJSON(message)
					if err != nil {
						log.Printf("error writing json to client: %v", err)
						// Don't delete here, let unregister handle it if conn breaks
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/kiel-live/kiel-live/pkg/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // Allow all connections
		return true
	},
}

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

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients     map[*websocket.Conn]map[string]struct{} // client -> set of subscribed topics
	broadcast   chan WebsocketMessage
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	subscribe   chan subscriptionRequest
	unsubscribe chan subscriptionRequest
	mu          sync.Mutex        // To protect clients map
	db          database.Database // Changed from store to db to match main.go
}

type subscriptionRequest struct {
	client *websocket.Conn
	topic  string
}

func newHub(db database.Database) *Hub {
	return &Hub{
		broadcast:   make(chan WebsocketMessage),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		subscribe:   make(chan subscriptionRequest),
		unsubscribe: make(chan subscriptionRequest),
		clients:     make(map[*websocket.Conn]map[string]struct{}),
		db:          db,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = make(map[string]struct{}) // Initialize empty set of topics
			h.mu.Unlock()
			log.Println("Client registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
		case req := <-h.subscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				clientTopics[req.topic] = struct{}{}
				log.Printf("Client subscribed to topic: %s", req.topic)
			} else {
				log.Printf("Failed to subscribe: client not registered: %v", req.client)
			}
			h.mu.Unlock()
		case req := <-h.unsubscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				delete(clientTopics, req.topic)
				log.Printf("Client unsubscribed from topic: %s", req.topic)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
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

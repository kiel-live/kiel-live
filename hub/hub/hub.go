package hub

import (
	"log"
	"sync"

	"github.com/kiel-live/kiel-live/pkg/database"

	"github.com/gorilla/websocket"
)

// websocketMessage defines the structure for messages sent over WebSocket.
type websocketMessage struct {
	Topic  string `json:"topic"`
	Action string `json:"action,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type subscriptionRequest struct {
	client *websocket.Conn
	topic  string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients     map[*websocket.Conn]map[string]struct{} // client -> set of subscribed topics
	mu          sync.Mutex                              // To protect clients map
	broadcast   chan websocketMessage
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	subscribe   chan subscriptionRequest
	unsubscribe chan subscriptionRequest
	db          database.Database
}

func NewHub(db database.Database) *Hub {
	return &Hub{
		broadcast:   make(chan websocketMessage),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		subscribe:   make(chan subscriptionRequest),
		unsubscribe: make(chan subscriptionRequest),
		clients:     make(map[*websocket.Conn]map[string]struct{}),
		db:          db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = make(map[string]struct{}) // Initialize empty set of topics
			h.mu.Unlock()
			h.sendStats()
			log.Println("Client registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
			h.sendStats()
		case req := <-h.subscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				clientTopics[req.topic] = struct{}{}
				log.Printf("Client subscribed to topic: %s", req.topic)
			} else {
				log.Printf("Failed to subscribe: client not registered: %v", req.client)
			}
			h.mu.Unlock()
			h.sendTopicList()
		case req := <-h.unsubscribe:
			h.mu.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				delete(clientTopics, req.topic)
				log.Printf("Client unsubscribed from topic: %s", req.topic)
			}
			h.mu.Unlock()
			h.sendTopicList()
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

func (h *Hub) RegisterClient(client *websocket.Conn) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *websocket.Conn) {
	h.unregister <- client
}

func (h *Hub) SubscribeClient(client *websocket.Conn, topic string) {
	h.subscribe <- subscriptionRequest{
		client: client,
		topic:  topic,
	}
}

func (h *Hub) UnsubscribeClient(client *websocket.Conn, topic string) {
	h.unsubscribe <- subscriptionRequest{
		client: client,
		topic:  topic,
	}
}

func (h *Hub) BroadcastMessage(topic string, action string, data any) {
	h.broadcast <- websocketMessage{
		Topic:  topic,
		Action: action,
		Data:   data,
	}
}

func (h *Hub) sendTopicList() {
	h.mu.Lock()
	defer h.mu.Unlock()

	topicsMap := map[string]struct{}{}
	for client := range h.clients {
		for topic := range h.clients[client] {
			topicsMap[topic] = struct{}{}
		}
	}

	topicsList := make([]string, 0, len(topicsMap))
	for topic := range topicsMap {
		topicsList = append(topicsList, topic)
	}

	h.BroadcastMessage("system.topics", "update", topicsList)
}

func (h *Hub) sendStats() {
	h.mu.Lock()
	defer h.mu.Unlock()

	stats := map[string]int{
		"clients": len(h.clients),
	}

	h.BroadcastMessage("system.stats", "update", stats)
}

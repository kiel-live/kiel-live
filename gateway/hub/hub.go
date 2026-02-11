package hub

import (
	"log"
	"sync"
	"time"

	"github.com/kiel-live/kiel-live/gateway/database"

	"github.com/gorilla/websocket"
)

// websocketMessage defines the structure for messages sent over WebSocket.
type websocketMessage struct {
	Topic  string    `json:"topic"`
	Action string    `json:"action,omitempty"`
	Data   any       `json:"data,omitempty"`
	SentAt time.Time `json:"sent_at"`
}

type subscriptionRequest struct {
	client *websocket.Conn
	topic  string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients      map[*websocket.Conn]map[string]struct{} // To protect clients map
	clientsMutex sync.RWMutex                            // Mutex to protect access to clients map
	broadcast    chan websocketMessage
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
	subscribe    chan subscriptionRequest
	unsubscribe  chan subscriptionRequest
	db           database.Database
}

func NewHub(db database.Database) *Hub {
	return &Hub{
		broadcast:   make(chan websocketMessage, 100),
		register:    make(chan *websocket.Conn, 10),
		unregister:  make(chan *websocket.Conn, 10),
		subscribe:   make(chan subscriptionRequest, 10),
		unsubscribe: make(chan subscriptionRequest, 10),
		clients:     make(map[*websocket.Conn]map[string]struct{}),
		db:          db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clientsMutex.Lock()
			h.clients[client] = make(map[string]struct{}) // Initialize empty set of topics
			h.clientsMutex.Unlock()
			log.Println("Client registered")
			h.sendStats()
		case client := <-h.unregister:
			h.clientsMutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Println("Client unregistered")
			}
			h.clientsMutex.Unlock()
			h.sendStats()
		case req := <-h.subscribe:
			h.clientsMutex.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				clientTopics[req.topic] = struct{}{}
				log.Printf("Client subscribed to topic: %s", req.topic)
			} else {
				log.Printf("Failed to subscribe: client not registered: %v", req.client)
			}
			h.clientsMutex.Unlock()
			h.sendTopicList()
		case req := <-h.unsubscribe:
			h.clientsMutex.Lock()
			if clientTopics, ok := h.clients[req.client]; ok {
				delete(clientTopics, req.topic)
				log.Printf("Client unsubscribed from topic: %s", req.topic)
			}
			h.clientsMutex.Unlock()
			h.sendTopicList()
		case message := <-h.broadcast:
			h.clientsMutex.RLock()
			for client, topics := range h.clients {
				if _, subscribed := topics[message.Topic]; subscribed {
					err := client.WriteJSON(message)
					if err != nil {
						log.Printf("error writing json to client: %v", err)
						// Don't delete here, let unregister handle it if conn breaks
					}
				}
			}
			h.clientsMutex.RUnlock()
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
		SentAt: time.Now(),
	}
}

func (h *Hub) sendTopicList() {
	h.clientsMutex.RLock()
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
	h.clientsMutex.RUnlock()

	h.BroadcastMessage("system.topics", "", topicsList)
}

func (h *Hub) sendStats() {
	h.clientsMutex.RLock()
	stats := map[string]int{
		"clients": len(h.clients),
	}
	h.clientsMutex.RUnlock()

	h.BroadcastMessage("system.stats", "", stats)
}

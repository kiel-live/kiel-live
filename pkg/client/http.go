package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/pkg/models"
)

type HTTPApiClient struct {
	baseURL          string
	httpClient       *http.Client
	token            string
	ws               *websocket.Conn
	subscribedTopics []string
	mu               sync.Mutex
	connected        bool

	connectionChangedHandler func(connected bool)
	topicsChangedHandler     func(topic string, subscribed bool)
}

func NewHTTPApiClient(baseURL string, token string) Client {
	return &HTTPApiClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		token:      token,
	}
}

func (h *HTTPApiClient) doRequest(method string, url string, request any, data any) error {
	reqBody := []byte{}

	if request != nil {
		raw, err := json.Marshal(request)
		if err != nil {
			return err
		}
		reqBody = raw
	}

	req, err := http.NewRequest(method, h.baseURL+url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if h.token != "" {
		req.Header.Set("Authorization", "Bearer "+h.token)
	}

	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}
	if data != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
			return err
		}
	}

	return nil
}

func (h *HTTPApiClient) get(url string, data any) error {
	return h.doRequest("GET", url, nil, data)
}

func (h *HTTPApiClient) put(url string, request any, data any) error {
	return h.doRequest("PUT", url, request, data)
}

func (h *HTTPApiClient) delete(url string, request any, data any) error {
	return h.doRequest("DELETE", url, request, data)
}

func (h *HTTPApiClient) Connect() error {
	u, err := url.Parse(h.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)
	u.Path = u.Path + "/ws"

	h.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		h.ws.Close()
		return err
	}

	// TODO: initially fetch subscribed topics
	// h.mu.Lock()
	// err = h.get("/topics", h.subscribedTopics)
	// h.mu.Unlock()
	// if err != nil {
	// 	h.ws.Close()
	// 	return fmt.Errorf("failed to get initial topics: %w", err)
	// }

	type SubscribeMessage struct {
		Topic  string `json:"topic"`
		Action string `json:"action"`
	}

	err = h.publish(SubscribeMessage{
		Topic:  "system.topics",
		Action: "subscribe",
	})
	if err != nil {
		h.ws.Close()
		return fmt.Errorf("failed to publish initial topics: %w", err)
	}

	h.connected = true
	if h.connectionChangedHandler != nil {
		h.connectionChangedHandler(true)
	}

	go func() {
		for {
			_, raw, err := h.ws.ReadMessage()
			if err != nil {
				return
			}

			var msg Message
			err = json.Unmarshal(raw, &msg)
			if err != nil {
				return
			}

			if msg.Topic == "system.topics" {
				var newTopics []string
				if err := json.Unmarshal([]byte(*msg.Data), &newTopics); err != nil {
					fmt.Printf("failed to unmarshal topics: %v\n", err)
					continue
				}

				h.mu.Lock()
				oldTopics := h.subscribedTopics
				h.subscribedTopics = newTopics
				h.mu.Unlock()

				if h.topicsChangedHandler != nil {
					for _, topic := range newTopics {
						if !slices.Contains(oldTopics, topic) {
							h.topicsChangedHandler(topic, true)
						}
					}

					for _, topic := range oldTopics {
						if !slices.Contains(newTopics, topic) {
							h.topicsChangedHandler(topic, false)
						}
					}
				}
			}
		}
	}()

	return nil
}

func (h *HTTPApiClient) publish(msg any) error {
	if h.ws == nil {
		return fmt.Errorf("not connected")
	}

	return h.ws.WriteJSON(msg)
}

func (h *HTTPApiClient) Disconnect() error {
	if h.ws != nil {
		err := h.ws.Close()
		h.ws = nil
		h.connected = false
		if h.connectionChangedHandler != nil {
			h.connectionChangedHandler(false)
		}
		return err
	}
	return nil
}

func (h *HTTPApiClient) IsConnected() bool {
	return h.ws != nil && h.connected
}

func (h *HTTPApiClient) SetOnConnectionChanged(handler func(connected bool)) {
	h.connectionChangedHandler = handler
}

func (h *HTTPApiClient) GetSubscribedTopics() []string {
	return h.subscribedTopics
}

func (h *HTTPApiClient) SetOnTopicsChanged(handler func(topic string, subscribed bool)) {
	h.topicsChangedHandler = handler
}

func (h *HTTPApiClient) UpdateStop(stop *models.Stop) error {
	return h.put("/stops/"+stop.ID, stop, nil)
}

func (h *HTTPApiClient) UpdateVehicle(vehicle *models.Vehicle) error {
	return h.put("/vehicles/"+vehicle.ID, vehicle, nil)
}

func (h *HTTPApiClient) UpdateTrip(trip *models.Trip) error {
	return h.put("/trips/"+trip.ID, trip, nil)
}

func (h *HTTPApiClient) DeleteStop(stopID string) error {
	return h.delete("/stops/"+stopID, nil, nil)
}

func (h *HTTPApiClient) DeleteVehicle(vehicleID string) error {
	return h.delete("/vehicles/"+vehicleID, nil, nil)
}

func (h *HTTPApiClient) DeleteTrip(tripID string) error {
	return h.delete("/trips/"+tripID, nil, nil)
}

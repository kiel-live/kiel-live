package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kiel-live/kiel-live/pkg/database"
	"github.com/kiel-live/kiel-live/pkg/models"

	"github.com/golang/geo/s2"
	"github.com/gorilla/websocket"
)

// APIServer holds dependencies for HTTP handlers
type APIServer struct {
	db  database.Database
	hub *Hub
	mux *http.ServeMux
}

// Helper methods for APIServer to register routes with HTTP methods
func (s *APIServer) GET(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("GET "+path, handler)
}

func (s *APIServer) POST(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("POST "+path, handler)
}

func (s *APIServer) PUT(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("PUT "+path, handler)
}

func (s *APIServer) DELETE(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("DELETE "+path, handler)
}

func (s *APIServer) Any(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc(path, handler)
}

func NewAPIServer(db database.Database, hub *Hub, mux *http.ServeMux) *APIServer {
	s := &APIServer{db: db, hub: hub, mux: mux}
	s.registerRoutes()
	return s
}

func (s *APIServer) registerRoutes() {
	// WebSocket
	s.Any("/ws", s.serveWs)

	// Vehicle CRUD
	s.GET("/vehicles", s.handleGetVehicles)
	s.GET("/vehicles/{id}", s.handleGetVehicle)
	s.PUT("/vehicles/{id}", s.handleUpdateVehicle)
	s.DELETE("/vehicles/{id}", s.handleDeleteVehicle)

	// Stop CRUD
	s.GET("/stops", s.handleGetStops)
	s.GET("/stops/{id}", s.handleGetStop)
	s.PUT("/stops/{id}", s.handleUpdateStop)
	s.DELETE("/stops/{id}", s.handleDeleteStop)

	// Trip CRUD
	s.GET("/trips/{id}", s.handleGetTrip)
	s.PUT("/trips/{id}", s.handleUpdateTrip)
	s.DELETE("/trips/{id}", s.handleDeleteTrip)

	// Route CRUD
	s.GET("/routes/{id}", s.handleGetRoute)
	s.PUT("/routes/{id}", s.handleUpdateRoute)
	s.DELETE("/routes/{id}", s.handleDeleteRoute)
}

func (s *APIServer) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	s.hub.register <- conn
	log.Println("WebSocket client connected")

	// Handle incoming messages from this client for subscriptions
	go func(c *websocket.Conn) {
		defer func() {
			s.hub.unregister <- c
			c.Close()
			log.Println("WebSocket client disconnected and unregistered")
		}()
		for {
			var msg ClientSubscriptionMessage
			if err := c.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error reading json from client: %v", err)
				}
				break // Exit loop on error
			}

			switch strings.ToLower(msg.Action) {
			case "subscribe":
				s.hub.subscribe <- subscriptionRequest{client: c, topic: msg.Topic}
			case "unsubscribe":
				s.hub.unsubscribe <- subscriptionRequest{client: c, topic: msg.Topic}
			default:
				log.Printf("Unknown action from client: %s", msg.Action)
			}
		}
	}(conn)
}

// Helper to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper to respond with error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Helper to convert s2.CellID to models.BoundingBox
// This is a simplified conversion; a more precise one might be needed depending on requirements.
func s2CellIDToBoundingBox(cellID s2.CellID) *models.BoundingBox {
	rect := s2.CellFromCellID(cellID).RectBound()
	return &models.BoundingBox{
		MinLat: rect.Lat.Lo,
		MinLng: rect.Lng.Lo,
		MaxLat: rect.Lat.Hi,
		MaxLng: rect.Lng.Hi,
	}
}

// --- Vehicle Handlers ---
func (s *APIServer) handleGetVehicles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cellIDStr := r.URL.Query().Get("s2CellID")
	if cellIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "s2CellID query parameter is required")
		return
	}
	cellID := s2.CellIDFromString(cellIDStr)
	if !cellID.IsValid() {
		respondWithError(w, http.StatusBadRequest, "Invalid s2CellID")
		return
	}

	bbox := s2CellIDToBoundingBox(cellID)
	vehicles, err := s.db.GetVehicles(ctx, &database.ListOptions{
		Location: bbox,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, vehicles)
}

func (s *APIServer) handleGetVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing vehicle ID")
		return
	}
	vehicle, err := s.db.GetVehicle(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if vehicle == nil {
		respondWithError(w, http.StatusNotFound, "Vehicle not found")
		return
	}
	respondWithJSON(w, http.StatusOK, vehicle)
}

func (s *APIServer) handleUpdateVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing vehicle ID")
		return
	}
	var v *models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if v.Location == nil {
		respondWithError(w, http.StatusBadRequest, "Vehicle location is required for update")
		return
	}
	if v.ID == "" {
		v.ID = id
	} else if v.ID != id {
		respondWithError(w, http.StatusBadRequest, "Vehicle ID in path and payload mismatch")
		return
	}

	oldVehicle, err := s.db.GetVehicle(ctx, id)
	var oldS2CellToken string
	if err == nil && oldVehicle != nil && oldVehicle.Location != nil {
		oldS2CellToken = oldVehicle.Location.GetCellID().ToToken()
	}

	err = s.db.SetVehicle(ctx, v)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Broadcast WebSocket updates
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("vehicles/%s", v.ID), Action: "updated", Data: v}
	if v.Location != nil {
		newS2CellToken := v.Location.GetCellID().ToToken()
		s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", newS2CellToken), Action: "updated", Data: v}
		if oldVehicle != nil && oldVehicle.Location != nil && oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Vehicle %s moved from S2 cell %s to %s", v.ID, oldS2CellToken, newS2CellToken)
			s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", oldS2CellToken), Action: "deleted", Data: map[string]string{"id": v.ID}}
		}
	}

	respondWithJSON(w, http.StatusOK, v)
}

func (s *APIServer) handleDeleteVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing vehicle ID")
		return
	}

	vehicle, err := s.db.GetVehicle(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching vehicle for delete: %v", err))
		return
	}
	if vehicle == nil {
		respondWithError(w, http.StatusNotFound, "Vehicle not found before delete attempt")
		return
	}

	if err := s.db.DeleteVehicle(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Broadcast WebSocket updates
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("vehicles/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	if vehicle.Location != nil {
		s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", vehicle.Location.GetCellID().ToToken()), Action: "deleted", Data: map[string]string{"id": id}}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// --- Stop Handlers ---
func (s *APIServer) handleGetStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cellIDStr := r.URL.Query().Get("s2CellID")
	if cellIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "s2CellID query parameter is required")
		return
	}
	cellID := s2.CellIDFromString(cellIDStr)
	if !cellID.IsValid() {
		respondWithError(w, http.StatusBadRequest, "Invalid s2CellID")
		return
	}

	bbox := s2CellIDToBoundingBox(cellID)
	stops, err := s.db.GetStops(ctx, &database.ListOptions{
		Location: bbox,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, stops)
}

func (s *APIServer) handleGetStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing stop ID")
		return
	}
	stop, err := s.db.GetStop(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if stop == nil {
		respondWithError(w, http.StatusNotFound, "Stop not found")
		return
	}
	respondWithJSON(w, http.StatusOK, stop)
}

func (s *APIServer) handleUpdateStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing stop ID")
		return
	}
	var stop *models.Stop
	if err := json.NewDecoder(r.Body).Decode(&stop); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if stop.Location == nil {
		respondWithError(w, http.StatusBadRequest, "Stop location is required for update")
		return
	}
	if stop.ID == "" {
		stop.ID = id
	} else if stop.ID != id {
		respondWithError(w, http.StatusBadRequest, "Stop ID in path and payload mismatch")
		return
	}

	oldStop, err := s.db.GetStop(ctx, id)
	var oldS2CellToken string
	if err == nil && oldStop != nil && oldStop.Location != nil {
		oldS2CellToken = oldStop.Location.GetCellID().ToToken()
	}

	err = s.db.SetStop(ctx, stop)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("stops/%s", stop.ID), Action: "updated", Data: stop}
	if stop.Location != nil {
		newS2CellToken := stop.Location.GetCellID().ToToken()
		s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", newS2CellToken), Action: "updated", Data: stop}
		if oldStop != nil && oldStop.Location != nil && oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Stop %s S2 cell changed from %s to %s", stop.ID, oldS2CellToken, newS2CellToken)
			s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", oldS2CellToken), Action: "deleted", Data: map[string]string{"id": stop.ID}}
		}
	}

	respondWithJSON(w, http.StatusOK, stop)
}

func (s *APIServer) handleDeleteStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing stop ID")
		return
	}

	stop, errDb := s.db.GetStop(ctx, id)
	if errDb != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching stop for delete: %v", errDb))
		return
	}
	if stop == nil {
		respondWithError(w, http.StatusNotFound, "Stop not found before delete attempt")
		return
	}

	if err := s.db.DeleteStop(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("stops/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	if stop.Location != nil {
		s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", stop.Location.GetCellID().ToToken()), Action: "deleted", Data: map[string]string{"id": id}}
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// --- Trip Handlers ---
func (s *APIServer) handleGetTrip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing trip ID")
		return
	}
	trip, err := s.db.GetTrip(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if trip == nil {
		respondWithError(w, http.StatusNotFound, "Trip not found")
		return
	}
	respondWithJSON(w, http.StatusOK, trip)
}

func (s *APIServer) handleUpdateTrip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing trip ID")
		return
	}
	var trip *models.Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if trip.ID == "" {
		trip.ID = id
	} else if trip.ID != id {
		respondWithError(w, http.StatusBadRequest, "Trip ID in path and payload mismatch")
		return
	}

	err := s.db.SetTrip(ctx, trip)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("trips/%s", trip.ID), Action: "updated", Data: trip}
	respondWithJSON(w, http.StatusOK, trip)
}

func (s *APIServer) handleDeleteTrip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing trip ID")
		return
	}
	if err := s.db.DeleteTrip(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("trips/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// --- Route Handlers ---
func (s *APIServer) handleGetRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}
	route, err := s.db.GetRoute(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if route == nil {
		respondWithError(w, http.StatusNotFound, "Route not found")
		return
	}
	respondWithJSON(w, http.StatusOK, route)
}

func (s *APIServer) handleUpdateRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}
	var route *models.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if route.ID == "" {
		route.ID = id
	} else if route.ID != id {
		respondWithError(w, http.StatusBadRequest, "Route ID in path and payload mismatch")
		return
	}

	err := s.db.SetRoute(ctx, route)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("routes/%s", route.ID), Action: "updated", Data: route}
	respondWithJSON(w, http.StatusOK, route)
}

func (s *APIServer) handleDeleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}
	if err := s.db.DeleteRoute(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast <- WebsocketMessage{Topic: fmt.Sprintf("routes/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

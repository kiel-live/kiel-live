package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/pkg/models"
)

// Helper to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
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

func (s *Server) broadcastItemUpdated(type2 string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(ItemTopic, type2, id), "update", data)
}

func (s *Server) broadcastMapItemUpdated(type2 string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(MapItemTopic, type2, id), "update", data)
}

func (s *Server) broadcastItemDeleted(type2 string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(ItemTopic, type2, id), "delete", data)
}

func (s *Server) broadcastMapItemDeleted(type2 string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(MapItemTopic, type2, id), "delete", data)
}

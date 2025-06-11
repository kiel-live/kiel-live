package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

func boundsFromQuery(values url.Values) (*models.BoundingBox, error) {
	getValue := func(key string) (float64, error) {
		if !values.Has(key) {
			return 0, fmt.Errorf("%s can't be empty", key)
		}
		value, err := strconv.ParseFloat(values.Get(key), 32)
		if err != nil {
			return 0, fmt.Errorf("%s is not a valid float: %w", key, err)
		}
		return value, nil
	}

	north, err := getValue("north")
	if err != nil {
		return nil, err
	}

	east, err := getValue("east")
	if err != nil {
		return nil, err
	}

	south, err := getValue("south")
	if err != nil {
		return nil, err
	}

	west, err := getValue("west")
	if err != nil {
		return nil, err
	}

	return &models.BoundingBox{
		North: north,
		East:  east,
		South: south,
		West:  west,
	}, nil
}

// Helper to convert s2.CellID to models.BoundingBox
// This is a simplified conversion; a more precise one might be needed depending on requirements.
func s2CellIDToBoundingBox(cellID s2.CellID) *models.BoundingBox {
	rect := s2.CellFromCellID(cellID).RectBound()
	return &models.BoundingBox{
		North: rect.Lat.Lo,
		East:  rect.Lng.Lo,
		South: rect.Lat.Hi,
		West:  rect.Lng.Hi,
	}
}

func (s *Server) broadcastItemUpdated(itemType string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(ItemTopic, itemType, id), "update", data)
}

func (s *Server) broadcastMapItemUpdated(itemType string, cellID string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(MapItemTopic, itemType, cellID), "update", data)
}

func (s *Server) broadcastItemDeleted(itemType string, id string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(ItemTopic, itemType, id), "delete", data)
}

func (s *Server) broadcastMapItemDeleted(itemType string, cellID string, data any) {
	s.hub.BroadcastMessage(fmt.Sprintf(MapItemTopic, itemType, cellID), "delete", data)
}

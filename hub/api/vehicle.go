package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/hub/hub"
	"github.com/kiel-live/kiel-live/pkg/database"
	"github.com/kiel-live/kiel-live/pkg/models"
)

// --- Vehicle Handlers ---
func (s *Server) handleGetVehicles(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleGetVehicle(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleUpdateVehicle(w http.ResponseWriter, r *http.Request) {
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
	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("vehicles/%s", v.ID), Action: "updated", Data: v}
	if v.Location != nil {
		newS2CellToken := v.Location.GetCellID().ToToken()
		s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", newS2CellToken), Action: "updated", Data: v}
		if oldVehicle != nil && oldVehicle.Location != nil && oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Vehicle %s moved from S2 cell %s to %s", v.ID, oldS2CellToken, newS2CellToken)
			s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", oldS2CellToken), Action: "deleted", Data: map[string]string{"id": v.ID}}
		}
	}

	respondWithJSON(w, http.StatusOK, v)
}

func (s *Server) handleDeleteVehicle(w http.ResponseWriter, r *http.Request) {
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
	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("vehicles/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	if vehicle.Location != nil {
		s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/vehicles", vehicle.Location.GetCellID().ToToken()), Action: "deleted", Data: map[string]string{"id": id}}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

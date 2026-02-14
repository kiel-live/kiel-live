package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kiel-live/kiel-live/gateway/database"
	"github.com/kiel-live/kiel-live/pkg/models"
)

func (s *Server) handleGetVehicles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bounds, err := boundsFromQuery(r.URL.Query())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	vehicles, err := s.db.GetVehicles(ctx, &database.ListOptions{
		Bounds: bounds,
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
	var vehicle *models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if vehicle.Location == nil {
		respondWithError(w, http.StatusBadRequest, "Vehicle location is required for update")
		return
	}
	if vehicle.ID == "" {
		vehicle.ID = id
	} else if vehicle.ID != id {
		respondWithError(w, http.StatusBadRequest, "Vehicle ID in path and payload mismatch")
		return
	}

	oldVehicle, err := s.db.GetVehicle(ctx, id)
	var oldS2CellToken string
	if err == nil && oldVehicle != nil && oldVehicle.Location != nil {
		oldS2CellToken = oldVehicle.Location.GetCellID().ToToken()
	}

	err = s.db.SetVehicle(ctx, vehicle)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.broadcastItemUpdated("vehicles", vehicle.ID, vehicle)
	if vehicle.Location != nil {
		newS2CellToken := vehicle.Location.GetCellID().ToToken()
		s.broadcastMapItemUpdated("vehicles", newS2CellToken, vehicle)
		if oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Vehicle %s moved from S2 cell %s to %s", vehicle.ID, oldS2CellToken, newS2CellToken)
			s.broadcastMapItemDeleted("vehicles", oldS2CellToken, vehicle)
		}
	}

	respondWithJSON(w, http.StatusOK, vehicle)
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

	s.broadcastItemDeleted("vehicles", id, vehicle)
	if vehicle.Location != nil {
		s.broadcastMapItemDeleted("vehicles", vehicle.Location.GetCellID().ToToken(), vehicle)
	}

	respondWithJSON(w, http.StatusOK, vehicle)
}

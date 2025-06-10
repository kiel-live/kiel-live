package api

import (
	"encoding/json"
	"net/http"

	"github.com/kiel-live/kiel-live/pkg/models"
)

func (s *Server) handleGetTrip(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleUpdateTrip(w http.ResponseWriter, r *http.Request) {
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
	s.broadcastItemUpdated("trips", trip.ID, trip)
	respondWithJSON(w, http.StatusOK, trip)
}

func (s *Server) handleDeleteTrip(w http.ResponseWriter, r *http.Request) {
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

	if err := s.db.DeleteTrip(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.broadcastMapItemDeleted("trips", id, trip)
	respondWithJSON(w, http.StatusOK, trip)
}

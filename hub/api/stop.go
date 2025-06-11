package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kiel-live/kiel-live/pkg/database"
	"github.com/kiel-live/kiel-live/pkg/models"
)

func (s *Server) handleGetStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bounds, err := boundsFromQuery(r.URL.Query())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	stops, err := s.db.GetStops(ctx, &database.ListOptions{
		Bounds: bounds,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, stops)
}

func (s *Server) handleGetStop(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleUpdateStop(w http.ResponseWriter, r *http.Request) {
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

	s.broadcastItemUpdated("stops", stop.ID, stop)
	if stop.Location != nil {
		newS2CellToken := stop.Location.GetCellID().ToToken()
		s.broadcastMapItemUpdated("stops", newS2CellToken, stop)
		if oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Stop %s S2 cell changed from %s to %s", stop.ID, oldS2CellToken, newS2CellToken)
			s.broadcastMapItemDeleted("stops", oldS2CellToken, stop)
		}
	}

	respondWithJSON(w, http.StatusOK, stop)
}

func (s *Server) handleDeleteStop(w http.ResponseWriter, r *http.Request) {
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

	s.broadcastItemDeleted("stops", id, stop)
	if stop.Location != nil {
		s.broadcastMapItemDeleted("stops", stop.Location.GetCellID().ToToken(), stop)
	}

	respondWithJSON(w, http.StatusOK, stop)
}

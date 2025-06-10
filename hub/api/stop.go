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

func (s *Server) handleGetStops(w http.ResponseWriter, r *http.Request) {
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

	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("stops/%s", stop.ID), Action: "updated", Data: stop}
	if stop.Location != nil {
		newS2CellToken := stop.Location.GetCellID().ToToken()
		s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", newS2CellToken), Action: "updated", Data: stop}
		if oldStop != nil && oldStop.Location != nil && oldS2CellToken != "" && oldS2CellToken != newS2CellToken {
			log.Printf("Stop %s S2 cell changed from %s to %s", stop.ID, oldS2CellToken, newS2CellToken)
			s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", oldS2CellToken), Action: "deleted", Data: map[string]string{"id": stop.ID}}
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
	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("stops/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	if stop.Location != nil {
		s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("s2cells/%s/stops", stop.Location.GetCellID().ToToken()), Action: "deleted", Data: map[string]string{"id": id}}
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

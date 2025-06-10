package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kiel-live/kiel-live/hub/hub"
	"github.com/kiel-live/kiel-live/pkg/models"
)

func (s *Server) handleGetRoute(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleUpdateRoute(w http.ResponseWriter, r *http.Request) {
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
	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("routes/%s", route.ID), Action: "updated", Data: route}
	respondWithJSON(w, http.StatusOK, route)
}

func (s *Server) handleDeleteRoute(w http.ResponseWriter, r *http.Request) {
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
	s.hub.Broadcast <- hub.WebsocketMessage{Topic: fmt.Sprintf("routes/%s", id), Action: "deleted", Data: map[string]string{"id": id}}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

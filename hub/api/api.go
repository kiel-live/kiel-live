package api

import (
	"net/http"

	"github.com/kiel-live/kiel-live/hub/hub"
	"github.com/kiel-live/kiel-live/pkg/database"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	db  database.Database
	hub *hub.Hub
	mux *http.ServeMux
}

func NewAPIServer(db database.Database, hub *hub.Hub, mux *http.ServeMux) *Server {
	s := &Server{db: db, hub: hub, mux: mux}
	s.registerRoutes()
	return s
}

// Helper methods for APIServer to register routes with HTTP methods
func (s *Server) GET(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("GET "+path, handler)
}

func (s *Server) POST(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("POST "+path, handler)
}

func (s *Server) PUT(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("PUT "+path, handler)
}

func (s *Server) DELETE(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("DELETE "+path, handler)
}

func (s *Server) Any(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc(path, handler)
}

func (s *Server) registerRoutes() {
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

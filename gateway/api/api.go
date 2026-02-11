package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kiel-live/kiel-live/gateway/database"
	"github.com/kiel-live/kiel-live/gateway/hub"
	"github.com/kiel-live/kiel-live/gateway/search"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	db     database.Database
	search search.Search
	hub    *hub.Hub
	mux    *http.ServeMux
	token  string
}

func NewAPIServer(db database.Database, search search.Search, hub *hub.Hub, mux *http.ServeMux, token string) *Server {
	s := &Server{
		db:     db,
		search: search,
		hub:    hub,
		mux:    mux,
		token:  token,
	}
	s.registerRoutes()
	return s
}

func corsWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			log.Printf("[%s] %s?%s", r.Method, r.URL.Path, r.URL.Query().Encode())
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler(w, r)
	}
}

// Helper methods for APIServer to register routes with HTTP methods
func (s *Server) GET(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("GET "+path, corsWrapper(handler))
}

func (s *Server) POST(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("POST "+path, corsWrapper(handler))
}

func (s *Server) PUT(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("PUT "+path, corsWrapper(handler))
}

func (s *Server) DELETE(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc("DELETE "+path, corsWrapper(handler))
}

func (s *Server) Any(path string, handler http.HandlerFunc) {
	s.mux.HandleFunc(path, corsWrapper(handler))
}

func (s *Server) WithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer "+s.token {
			code := http.StatusUnauthorized
			http.Error(w, fmt.Sprintf("%d %s", code, http.StatusText(code)), code)
			return
		}

		handler(w, r)
	}
}

func (s *Server) registerRoutes() {
	// WebSocket
	s.Any("/ws", s.serveWs)

	// Vehicle CRUD
	s.GET("/vehicles", s.handleGetVehicles)
	s.GET("/vehicles/{id}", s.handleGetVehicle)
	s.PUT("/vehicles/{id}", s.WithAuth(s.handleUpdateVehicle))
	s.DELETE("/vehicles/{id}", s.WithAuth(s.handleDeleteVehicle))

	// Stop CRUD
	s.GET("/stops", s.handleGetStops)
	s.GET("/stops/{id}", s.handleGetStop)
	s.PUT("/stops/{id}", s.WithAuth(s.handleUpdateStop))
	s.DELETE("/stops/{id}", s.WithAuth(s.handleDeleteStop))

	// Trip CRUD
	s.GET("/trips/{id}", s.handleGetTrip)
	s.PUT("/trips/{id}", s.WithAuth(s.handleUpdateTrip))
	s.DELETE("/trips/{id}", s.WithAuth(s.handleDeleteTrip))
}

package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/pkg/metrics"
	"github.com/kiel-live/kiel-live/hub/store"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	store          *store.Store
	subMgr         *subscriptionManager
	token          string
	reg            *metrics.Registry
	clientCount    atomic.Int64
	collectorCount atomic.Int64
}

func New(st *store.Store, token string) *Server {
	reg := metrics.New()
	srv := &Server{
		store:  st,
		subMgr: newSubscriptionManager(),
		token:  token,
		reg:    reg,
	}
	st.RegisterMetrics(reg)
	reg.Register("clients", func() any { return srv.clientCount.Load() })
	reg.Register("collectors", func() any { return srv.collectorCount.Load() })
	srv.subMgr.registerMetrics(reg)
	return srv
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/client", s.handleClientWS)
	mux.HandleFunc("/ws/collector", s.handleCollectorWS)
	mux.HandleFunc("/api/status", s.handleStatus)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.reg.Snapshot())
}

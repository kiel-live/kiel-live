package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/metrics"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Default WebSocket keepalive timing for both /ws/client and /ws/collector.
// pongWait should be a small multiple of pingInterval so a couple of missed
// pings are tolerated before a connection is considered dead.
const (
	defaultPingInterval = 30 * time.Second
	defaultPongWait     = 60 * time.Second
)

type Server struct {
	store          *store.Store
	subMgr         *subscriptionManager
	token          string
	reg            *metrics.Registry
	clientCount    atomic.Int64
	collectorCount atomic.Int64
	pingInterval   time.Duration
	pongWait       time.Duration
}

func New(st *store.Store, token string) *Server {
	reg := metrics.New()
	srv := &Server{
		store:        st,
		subMgr:       newSubscriptionManager(),
		token:        token,
		reg:          reg,
		pingInterval: defaultPingInterval,
		pongWait:     defaultPongWait,
	}
	st.RegisterMetrics(reg)
	reg.Register("clients", func() any { return srv.clientCount.Load() })
	reg.Register("collectors", func() any { return srv.collectorCount.Load() })
	srv.subMgr.registerMetrics(reg)
	return srv
}

// SetKeepalive overrides the WebSocket ping interval and pong wait timeout
// used to detect dead /ws/client and /ws/collector connections. Intended for
// tests that want a fast keepalive cycle; production code should rely on the
// defaults set in New.
func (s *Server) SetKeepalive(pingInterval, pongWait time.Duration) {
	s.pingInterval = pingInterval
	s.pongWait = pongWait
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

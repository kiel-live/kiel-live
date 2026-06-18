package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/protocol"
)

type collectorConn struct {
	ws     *websocket.Conn
	st     *store.Store
	subMgr *subscriptionManager
	log    *slog.Logger

	writeMu sync.Mutex
}

func (s *Server) handleCollectorWS(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	expected := "Bearer " + s.token
	if auth != expected {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("collector ws upgrade failed", "error", err)
		return
	}

	c := &collectorConn{
		ws:     ws,
		st:     s.store,
		subMgr: s.subMgr,
		log:    slog.With("role", "collector", "remote", r.RemoteAddr),
	}
	s.collectorCount.Add(1)
	c.log.Info("collector connected")
	c.run()
	s.collectorCount.Add(-1)
	c.log.Info("collector disconnected")
}

func (c *collectorConn) run() {
	defer c.cleanup()

	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log.Debug("collector read error", "error", err)
			}
			return
		}

		var env protocol.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			c.log.Error("parse error", "error", err)
			continue
		}

		c.handleEnvelope(env)
	}
}

func (c *collectorConn) handleEnvelope(env protocol.Envelope) {
	switch env.Method {
	case protocol.MethodSubscribe:
		if env.Topic == protocol.TopicSubscriptions {
			c.subMgr.registerCollector(c)
		}

	case protocol.MethodUpdate:
		_, kind, _ := protocol.ParseTopic(env.Topic)
		if env.Data == nil {
			c.log.Warn("update without data", "topic", env.Topic)
			return
		}
		switch kind {
		case "stop":
			var stop models.Stop
			if err := json.Unmarshal(*env.Data, &stop); err != nil {
				c.log.Error("unmarshal stop failed", "error", err)
				return
			}
			c.st.UpsertStop(&stop)
		case "vehicle":
			var vehicle models.Vehicle
			if err := json.Unmarshal(*env.Data, &vehicle); err != nil {
				c.log.Error("unmarshal vehicle failed", "error", err)
				return
			}
			c.st.UpsertVehicle(&vehicle)
		case "trip":
			var trip models.Trip
			if err := json.Unmarshal(*env.Data, &trip); err != nil {
				c.log.Error("unmarshal trip failed", "error", err)
				return
			}
			c.st.UpsertTrip(&trip)
		}

	case protocol.MethodDelete:
		_, kind, id := protocol.ParseTopic(env.Topic)
		switch kind {
		case "stop":
			c.st.DeleteStop(id)
		case "vehicle":
			c.st.DeleteVehicle(id)
		case "trip":
			c.st.DeleteTrip(id)
		}
	}
}

func (c *collectorConn) cleanup() {
	c.subMgr.unregisterCollector(c)
	c.ws.Close()
}

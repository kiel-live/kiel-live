package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang/geo/s2"
	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/protocol"
)

const eventBufSize = 512

type clientConn struct {
	ws     *websocket.Conn
	st     *store.Store
	subMgr *subscriptionManager
	log    *slog.Logger

	// state — only accessed in the main event loop goroutine
	mapKindSubs map[string]bool            // "stop" or "vehicle" → opted in
	cellSubs    map[string]map[string]bool // kind → set of cellToken strings
	detailSubs  map[string]bool            // hub topic → subscribed
	lastVP      *protocol.ViewportData     // last viewport, replayed when map sub arrives late

	inboundCh chan protocol.Envelope
	eventCh   chan store.Event
	done      chan struct{}
}

func (s *Server) handleClientWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("client ws upgrade failed", "error", err)
		return
	}

	c := &clientConn{
		ws:          ws,
		st:          s.store,
		subMgr:      s.subMgr,
		log:         slog.With("role", "client", "remote", r.RemoteAddr),
		mapKindSubs: make(map[string]bool),
		cellSubs:    map[string]map[string]bool{"stop": {}, "vehicle": {}},
		detailSubs:  make(map[string]bool),
		inboundCh:   make(chan protocol.Envelope, 32),
		eventCh:     make(chan store.Event, eventBufSize),
		done:        make(chan struct{}),
	}
	s.clientCount.Add(1)
	c.log.Info("client connected")
	startKeepalive(ws, s.pingInterval, s.pongWait, c.done)
	c.run()
	s.clientCount.Add(-1)
	c.log.Info("client disconnected")
}

func (c *clientConn) run() {
	defer func() {
		close(c.done)
		// unsubscribe all pub/sub topics
		for kind, tokens := range c.cellSubs {
			for token := range tokens {
				c.st.PubSub.Unsubscribe("map."+kind+"."+token, c.eventCh)
			}
		}
		for topic := range c.detailSubs {
			c.st.PubSub.Unsubscribe(topic, c.eventCh)
			c.subMgr.unsubscribe(topic)
		}
		c.ws.Close()
	}()

	go c.readWS()

	for {
		select {
		case env, ok := <-c.inboundCh:
			if !ok {
				return
			}
			c.handleMessage(env)
		case event := <-c.eventCh:
			c.handlePubSubEvent(event)
		}
	}
}

func (c *clientConn) readWS() {
	defer close(c.inboundCh)
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log.Debug("client read error", "error", err)
			}
			return
		}
		var env protocol.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			c.log.Error("client parse error", "error", err)
			continue
		}
		select {
		case c.inboundCh <- env:
		case <-c.done:
			return
		}
	}
}

func (c *clientConn) handleMessage(env protocol.Envelope) {
	switch env.Method {
	case protocol.MethodViewport:
		if env.Data == nil {
			return
		}
		var vp protocol.ViewportData
		if err := json.Unmarshal(*env.Data, &vp); err != nil {
			c.log.Error("viewport parse error", "error", err)
			return
		}
		c.lastVP = &vp
		c.applyViewport(vp)

	case protocol.MethodSubscribe:
		topic := env.Topic
		if topic == "" {
			return
		}
		// "map.stop" or "map.vehicle" — opt-in to map kind streaming
		if topic == "map.stop" || topic == "map.vehicle" {
			kind := topic[4:] // "stop" or "vehicle"
			c.subscribeMapKind(kind)
			return
		}
		// detail subscription: "stop.{id}", "vehicle.{id}", "trip.{id}"
		if !c.detailSubs[topic] {
			c.detailSubs[topic] = true
			c.st.PubSub.Subscribe(topic, c.eventCh)
			c.subMgr.subscribe(topic)
			// send current entity immediately if available
			_, kind, id := protocol.ParseTopic(topic)
			switch kind {
			case "stop":
				if s := c.st.GetStop(id); s != nil {
					c.sendJSON(protocol.Envelope{Method: protocol.MethodUpdate, Topic: topic, Data: toRaw(s)})
				}
			case "vehicle":
				if v := c.st.GetVehicle(id); v != nil {
					c.sendJSON(protocol.Envelope{Method: protocol.MethodUpdate, Topic: topic, Data: toRaw(v)})
				}
			case "trip":
				if t := c.st.GetTrip(id); t != nil {
					c.sendJSON(protocol.Envelope{Method: protocol.MethodUpdate, Topic: topic, Data: toRaw(t)})
				}
			}
		}

	case protocol.MethodUnsubscribe:
		topic := env.Topic
		if topic == "map.stop" || topic == "map.vehicle" {
			kind := topic[4:]
			c.unsubscribeMapKind(kind)
			return
		}
		if c.detailSubs[topic] {
			delete(c.detailSubs, topic)
			c.st.PubSub.Unsubscribe(topic, c.eventCh)
			c.subMgr.unsubscribe(topic)
		}

	case protocol.MethodRequest:
		if env.Topic == protocol.TopicSearch {
			c.handleSearch(env)
		}
	}
}

// subscribeMapKind opts the client into map streaming for the given kind ("stop" or "vehicle").
// If a viewport has already been received, cells are subscribed immediately.
func (c *clientConn) subscribeMapKind(kind string) {
	if c.mapKindSubs[kind] {
		return
	}
	c.mapKindSubs[kind] = true
	if c.lastVP != nil {
		c.syncCells(kind, *c.lastVP)
	}
}

// unsubscribeMapKind opts the client out and clears all cell subscriptions for kind.
func (c *clientConn) unsubscribeMapKind(kind string) {
	if !c.mapKindSubs[kind] {
		return
	}
	delete(c.mapKindSubs, kind)
	c.clearCells(kind)
}

// applyViewport computes the covering cells and syncs subscriptions for all opted-in kinds.
func (c *clientConn) applyViewport(vp protocol.ViewportData) {
	for kind := range c.mapKindSubs {
		c.syncCells(kind, vp)
	}
}

// syncCells subscribes to newly covered cells and unsubscribes from departed ones,
// sending a snapshot for each newly subscribed cell.
func (c *clientConn) syncCells(kind string, vp protocol.ViewportData) {
	bbox := models.BoundingBox{North: vp.North, East: vp.East, South: vp.South, West: vp.West}
	newCells := bbox.GetCellIDs()

	newTokens := make(map[string]s2.CellID, len(newCells))
	for _, cell := range newCells {
		newTokens[cell.ToToken()] = cell
	}

	current := c.cellSubs[kind]

	// Subscribe to newly entered cells
	for token, cell := range newTokens {
		if current[token] {
			continue
		}
		pubTopic := "map." + kind + "." + token
		// Subscribe before snapshot to avoid missing concurrent updates
		c.st.PubSub.Subscribe(pubTopic, c.eventCh)
		current[token] = true
		c.sendCellSnapshot(kind, cell)
	}

	// Unsubscribe from departed cells and send deletes
	for token := range current {
		if newTokens[token] != 0 {
			continue
		}
		pubTopic := "map." + kind + "." + token
		c.st.PubSub.Unsubscribe(pubTopic, c.eventCh)
		delete(current, token)
		c.sendCellDeletes(kind, s2.CellIDFromToken(token))
	}
}

// clearCells unsubscribes all cells for kind and sends delete messages.
func (c *clientConn) clearCells(kind string) {
	for token := range c.cellSubs[kind] {
		pubTopic := "map." + kind + "." + token
		c.st.PubSub.Unsubscribe(pubTopic, c.eventCh)
		c.sendCellDeletes(kind, s2.CellIDFromToken(token))
	}
	c.cellSubs[kind] = make(map[string]bool)
}

// sendCellSnapshot sends a snapshot of all current markers in cell to the client.
func (c *clientConn) sendCellSnapshot(kind string, cell s2.CellID) {
	cells := []s2.CellID{cell}
	var items []any
	switch kind {
	case "stop":
		stops, _ := c.st.GetMarkersInCells(cells)
		items = stopsToAny(stops)
	case "vehicle":
		_, vehicles := c.st.GetMarkersInCells(cells)
		items = vehiclesToAny(vehicles)
	}
	c.sendSnapshot(kind, items)
}

// sendCellDeletes sends delete messages for all current entities in cell.
func (c *clientConn) sendCellDeletes(kind string, cell s2.CellID) {
	cells := []s2.CellID{cell}
	switch kind {
	case "stop":
		stops, _ := c.st.GetMarkersInCells(cells)
		for _, s := range stops {
			c.sendJSON(protocol.Envelope{Method: protocol.MethodDelete, Topic: protocol.FormatMapTopic(kind, s.ID)})
		}
	case "vehicle":
		_, vehicles := c.st.GetMarkersInCells(cells)
		for _, v := range vehicles {
			c.sendJSON(protocol.Envelope{Method: protocol.MethodDelete, Topic: protocol.FormatMapTopic(kind, v.ID)})
		}
	}
}

func (c *clientConn) handlePubSubEvent(e store.Event) {
	if strings.HasPrefix(e.Topic, "map.") {
		// Map event: send sparse marker to client
		wireTopic := protocol.FormatMapTopic(e.Kind, e.ID)
		if e.Action == store.EventActionUpdate {
			if e.Marker == nil {
				return
			}
			c.sendJSON(protocol.Envelope{Method: protocol.MethodUpdate, Topic: wireTopic, Data: toRaw(e.Marker)})
		} else {
			c.sendJSON(protocol.Envelope{Method: protocol.MethodDelete, Topic: wireTopic})
		}
		return
	}
	// Detail event: send full entity
	if e.Action == store.EventActionUpdate {
		if e.Entity == nil {
			return
		}
		c.sendJSON(protocol.Envelope{Method: protocol.MethodUpdate, Topic: e.Topic, Data: toRaw(e.Entity)})
	}
}

const defaultSearchLimit = 20
const maxSearchLimit = 100

func (c *clientConn) handleSearch(env protocol.Envelope) {
	if env.Data == nil {
		c.sendError(env.RID, protocol.TopicSearch, "missing data")
		return
	}
	var req protocol.SearchRequest
	if err := json.Unmarshal(*env.Data, &req); err != nil {
		c.sendError(env.RID, protocol.TopicSearch, "invalid request")
		return
	}
	if req.Query == "" {
		c.sendError(env.RID, protocol.TopicSearch, "query is required")
		return
	}
	limit := req.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	} else if limit > maxSearchLimit {
		limit = maxSearchLimit
	}

	var bounds *models.BoundingBox
	if b := req.Bounds; b != nil {
		bounds = &models.BoundingBox{North: b.North, East: b.East, South: b.South, West: b.West}
	}
	results := c.st.Search(req.Query, bounds, limit)
	data, err := json.Marshal(results)
	if err != nil {
		c.sendError(env.RID, protocol.TopicSearch, "internal error")
		return
	}
	raw := json.RawMessage(data)
	c.sendJSON(protocol.Envelope{
		RID:    env.RID,
		Method: protocol.MethodResult,
		Topic:  protocol.TopicSearch,
		Data:   &raw,
	})
}

func (c *clientConn) sendError(rid, topic, msg string) {
	data, _ := json.Marshal(msg)
	raw := json.RawMessage(data)
	c.sendJSON(protocol.Envelope{
		RID:    rid,
		Method: protocol.MethodError,
		Topic:  topic,
		Data:   &raw,
	})
}

func (c *clientConn) sendSnapshot(kind string, items []any) {
	data, err := json.Marshal(items)
	if err != nil {
		c.log.Error("snapshot marshal error", "error", err)
		return
	}
	raw := json.RawMessage(data)
	c.sendJSON(protocol.Envelope{
		Method: protocol.MethodSnapshot,
		Topic:  protocol.FormatMapKindTopic(kind),
		Data:   &raw,
	})
}

func (c *clientConn) sendJSON(env protocol.Envelope) {
	data, err := json.Marshal(env)
	if err != nil {
		c.log.Error("marshal error", "error", err)
		return
	}
	if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		c.log.Debug("write error", "error", err)
	}
}

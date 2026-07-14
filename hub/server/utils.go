package server

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/pkg/models"
)

// keepaliveWriteWait bounds how long a single ping control frame write may block.
// It is intentionally not configurable — it only guards the write call itself,
// not the ping cadence or dead-peer detection window.
const keepaliveWriteWait = 5 * time.Second

// startKeepalive sends periodic WS pings on ws and resets ws's read deadline
// whenever a pong arrives. If no pong (or any other message) arrives within
// pongWait, the next ReadMessage call on ws returns an error, so the caller's
// existing read-loop error handling tears the connection down — this only
// adds a way to detect a half-open connection, not new cleanup logic.
//
// The ping goroutine exits when done is closed, which callers must do as part
// of their existing connection teardown.
func startKeepalive(ws *websocket.Conn, pingInterval, pongWait time.Duration, done <-chan struct{}) {
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		return ws.SetReadDeadline(time.Now().Add(pongWait))
	})

	go func() {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(keepaliveWriteWait)); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()
}

func toRaw(v any) *json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	raw := json.RawMessage(data)
	return &raw
}

func stopsToAny(stops []*models.Stop) []any {
	out := make([]any, len(stops))
	for i, s := range stops {
		out[i] = s
	}
	return out
}

func vehiclesToAny(vehicles []*models.Vehicle) []any {
	out := make([]any, len(vehicles))
	for i, v := range vehicles {
		out[i] = v
	}
	return out
}

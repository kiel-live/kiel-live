package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

const (
	IDPrefix   = "ais-"
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

var wsWriteMutex sync.Mutex
var keepAliveCancelMu sync.Mutex
var keepAliveCancel context.CancelFunc

type shipPosition struct {
	lat float64
	lon float64
}

var lastPositions = make(map[int]shipPosition)
var lastPositionsMu sync.Mutex

func wsConnect() (*websocket.Conn, error) {
	url := "wss://stream.aisstream.io/v0/stream"
	ws, http, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	_ = http

	return ws, nil
}

func configureConnection(ws *websocket.Conn) {
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// cancel any previous keep-alive goroutine
	keepAliveCancelMu.Lock()
	if keepAliveCancel != nil {
		keepAliveCancel()
		keepAliveCancel = nil
	}
	var ctx context.Context
	ctx, keepAliveCancel = context.WithCancel(context.Background())
	keepAliveCancelMu.Unlock()

	go startKeepAlive(ctx, ws)
}

func startKeepAlive(ctx context.Context, ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			wsWriteMutex.Lock()
			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				wsWriteMutex.Unlock()
				slog.Error("Error setting write deadline", "error", err)
				return
			}
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				wsWriteMutex.Unlock()
				return
			}
			wsWriteMutex.Unlock()
		}
	}
}

func subscribeToStream(ws *websocket.Conn, subMsg aisstream.SubscriptionMessage) error {
	subMsgBytes, err := json.Marshal(subMsg)
	if err != nil {
		return err
	}

	wsWriteMutex.Lock()
	defer wsWriteMutex.Unlock()
	return ws.WriteMessage(websocket.TextMessage, subMsgBytes)
}

func main() {
	collector.New(collector.Options{
		Name:    "🚢 AIS",
		Execute: run,
	}).Run()
}

func run(_ context.Context, c client.Client) error {
	apiKey := os.Getenv("AISSTREAM_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("please provide an API key for the AIS stream with AISSTREAM_API_KEY")
	}

	ws, err := wsConnect()
	if err != nil {
		return err
	}
	slog.Info("Connected to AIS stream")
	defer ws.Close()

	configureConnection(ws)

	subMsg := aisstream.SubscriptionMessage{
		APIKey:        apiKey,
		BoundingBoxes: [][][]float64{{{54.0, 10.0}, {55.0, 11.0}}},
		FiltersShipMMSI: []string{
			"211865680", // Wik
			"211341930", // Gaarden
			"211848130", // Friedrichsort
			"211872380", // Wellingdorf
			"218035310", // Laboe
			"218039370", // Dietrichsdorf
			"211399920", // Duesternbrook
			"211549030", // Adler 1
		},
	}

	if err := subscribeToStream(ws, subMsg); err != nil {
		return err
	}
	slog.Info("Subscribed to AIS stream")

	for {
		_, p, err := ws.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			slog.Warn("Websocket connection closed, reconnecting...")
			ws, err = wsConnect()
			if err != nil {
				return err
			}
			slog.Info("Reconnected to AIS stream")
			configureConnection(ws)
			if err = subscribeToStream(ws, subMsg); err != nil {
				return err
			}
			slog.Info("Resubscribed to AIS stream")
			continue
		} else if err != nil {
			return err
		}

		var packet aisstream.AisStreamMessage
		if err = json.Unmarshal(p, &packet); err != nil {
			slog.Error("failed to decode AIS packet", "error", err)
			continue
		}

		var shipName string
		// field may or may not be populated
		if packetShipName, ok := packet.MetaData["ShipName"]; ok {
			shipName = packetShipName.(string)
		}

		switch packet.MessageType {
		case aisstream.POSITION_REPORT:
			var positionReport aisstream.PositionReport
			positionReport = *packet.Message.PositionReport
			mmsi := int(positionReport.UserID)
			slog.Debug("position report", "mmsi", mmsi, "ship_name", shipName, "latitude", positionReport.Latitude, "longitude", positionReport.Longitude, "true_heading", positionReport.TrueHeading)
			location := &models.Location{
				Longitude: int(positionReport.Longitude * 3600000),
				Latitude:  int(positionReport.Latitude * 3600000),
			}
			if positionReport.TrueHeading != 511 { // 511 means "not available"
				heading := int(positionReport.TrueHeading)
				location.Heading = &heading
			} else {
				lastPositionsMu.Lock()
				previousPosition, ok := lastPositions[mmsi]
				lastPositionsMu.Unlock()
				if ok && calculateDistanceMeters(previousPosition.lat, previousPosition.lon, positionReport.Latitude, positionReport.Longitude) > 10 {
					heading := calculateBearing(previousPosition.lat, previousPosition.lon, positionReport.Latitude, positionReport.Longitude)
					location.Heading = &heading
				}
			}
			lastPositionsMu.Lock()
			lastPositions[mmsi] = shipPosition{lat: positionReport.Latitude, lon: positionReport.Longitude}
			lastPositionsMu.Unlock()

			vehicle := &models.Vehicle{
				ID:          IDPrefix + fmt.Sprint(positionReport.UserID),
				Provider:    "ais",
				Name:        shipName,
				Description: fmt.Sprintf("Die Live-Position der Fähre \"%s\".\n\nDie Position wird regelmäßig über AIS aktualisiert.", shipName),
				Type:        models.VehicleTypeFerry,
				State:       "onfire", // TODO
				Location:    location,
			}

			if err = c.UpdateVehicle(vehicle); err != nil {
				slog.Error(err.Error())
				continue
			}
		}
	}
}

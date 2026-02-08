package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
	log "github.com/sirupsen/logrus"
)

const IDPrefix = "ais-"

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

var wsWriteMutex sync.Mutex

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
	go startKeepAlive(ws)
}

func startKeepAlive(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for range ticker.C {
		wsWriteMutex.Lock()
		if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			wsWriteMutex.Unlock()
			log.Errorln("Error setting write deadline:", err)
			return
		}
		if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			wsWriteMutex.Unlock()
			return
		}
		wsWriteMutex.Unlock()
	}
}

func subscribeToStream(ws *websocket.Conn, subMsg aisstream.SubscriptionMessage) error {
	subMsgBytes, err := json.Marshal(subMsg)
	if err != nil {
		return err
	}

	wsWriteMutex.Lock()
	defer wsWriteMutex.Unlock()
	if err := ws.WriteMessage(websocket.TextMessage, subMsgBytes); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Infof("Kiel-Live AIS collector version %s", "1.0.0") // TODO use proper version

	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	apiKey := os.Getenv("AISSTREAM_API_KEY")
	if apiKey == "" {
		log.Fatalln("Please provide an API key for the AIS stream with AISSTREAM_API_KEY")
	}

	server := os.Getenv("COLLECTOR_SERVER")
	if server == "" {
		log.Fatalln("Please provide a server address for the collector with COLLECTOR_SERVER")
	}

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the collector with MANAGER_TOKEN")
	}

	c := client.NewClient(server, token)
	err = c.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer func() {
		err := c.Disconnect()
		if err != nil {
			log.Error(err)
		}
	}()

	ws, err := wsConnect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Infoln("Connected to AIS stream")
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

	err = subscribeToStream(ws, subMsg)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Infoln("Subscribed to AIS stream")

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				// reconnect
				log.Warnln("Websocket connection closed, reconnecting...")
				ws, err = wsConnect()
				if err != nil {
					log.Fatalln(err)
				}
				log.Infoln("Reconnected to AIS stream")
				configureConnection(ws)
				err = subscribeToStream(ws, subMsg)
				if err != nil {
					log.Fatalln(err)
				}
				log.Infoln("Resubscribed to AIS stream")
				continue
			}
			log.Fatalln(err)
		}
		var packet aisstream.AisStreamMessage

		err = json.Unmarshal(p, &packet)
		if err != nil {
			log.Errorln(err)
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
			log.Debugf("MMSI: %d Ship Name: %s Latitude: %f Longitude: %f",
				positionReport.UserID, shipName, positionReport.Latitude, positionReport.Longitude)
			vehicle := &models.Vehicle{
				ID:          IDPrefix + fmt.Sprint(positionReport.UserID),
				Provider:    "ais",
				Name:        shipName,
				Description: fmt.Sprintf("Die Live-Position der Fähre \"%s\".\n\nDie Position wird regelmäßig über AIS aktualisiert.", shipName),
				Type:        models.VehicleTypeFerry,
				State:       "onfire", // TODO
				Location: &models.Location{
					Longitude: int(positionReport.Longitude * 3600000),
					Latitude:  int(positionReport.Latitude * 3600000),
					Heading:   int(positionReport.TrueHeading),
				},
			}

			err = c.UpdateVehicle(vehicle)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

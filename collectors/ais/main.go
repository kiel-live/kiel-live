package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	ais "github.com/BertoldVdb/go-ais"
	"github.com/BertoldVdb/go-ais/aisnmea"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

const IDPrefix = "ais-"

func main() {
	log.Infof("Kiel-Live AIS collector version %s", "1.0.0") // TODO use proper version

	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	address := os.Getenv("UDP_ADDRESS")
	if address == "" {
		log.Fatalln("Please provide a UDP address for the collector with UDP_ADDRESS")
	}

	port := os.Getenv("UDP_PORT")
	if port == "" {
		log.Fatalln("Please provide a UDP port for the collector with UDP_PORT")
	}

	server := os.Getenv("COLLECTOR_SERVER")
	if server == "" {
		log.Fatalln("Please provide a server address for the collector with COLLECTOR_SERVER")
	}

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the collector with MANAGER_TOKEN")
	}

	s, err := net.ResolveUDPAddr("udp4", address+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	c := client.NewNatsClient(server, client.WithAuth("collector", token))
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

	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Error(err)
			continue
		}

		nm := aisnmea.NMEACodecNew(ais.CodecNew(false, false))

		decoded, err := nm.ParseSentence(string(buffer[0 : n-1]))
		if err != nil {
			log.Error(err)
			continue
		}
		if decoded != nil {
			if decoded.Packet.GetHeader().MessageID == 1 || decoded.Packet.GetHeader().MessageID == 2 || decoded.Packet.GetHeader().MessageID == 3 {
				positionReportPacket := decoded.Packet.(ais.PositionReport)

				vehicle := protocol.Vehicle{
					ID:       IDPrefix + fmt.Sprint(positionReportPacket.UserID),
					Provider: "ais",
					Type:     protocol.VehicleTypeFerry,
					State:    "onfire", // TODO
					Location: protocol.Location{
						Heading:   int(positionReportPacket.TrueHeading),
						Longitude: int(positionReportPacket.Longitude * 3600000),
						Latitude:  int(positionReportPacket.Latitude * 3600000),
					},
				}

				subject := fmt.Sprintf(protocol.TopicMapVehicle, vehicle.ID)

				jsonData, err := json.Marshal(vehicle)
				if err != nil {
					log.Error(err)
					continue
				}

				err = c.Publish(subject, string(jsonData))
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("🚴‍♀️ Nextbike collector version %s", "1.0.0") // TODO use proper version

	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	if os.Getenv("LOG") == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	server := os.Getenv("COLLECTOR_SERVER")
	if server == "" {
		log.Fatalln("Please provide a server address for the collector with COLLECTOR_SERVER")
	}

	log.Println("🚀 Connecting to server...", server)

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the collector with COLLECTOR_TOKEN")
	}

	cityIds := os.Getenv("NEXT_BIKE_CITY_IDS")
	if token == "" {
		log.Fatalln("Please provide a comma separated list of next-bike city ids with NEXT_BIKE_CITY_IDS (exp: '613,195' for Kiel & Mannheim)")
	}

	c := client.NewClient(server, client.WithAuth("collector", token))
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

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err = s.Every(5).Seconds().Do(func() error {
		if !c.IsConnected() {
			return nil
		}

		resp, err := http.Get("https://api.nextbike.net/maps/nextbike-live.json?city=" + cityIds)
		if err != nil {
			return err
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		nextbikeResp := NextbikeResponse{}
		err = json.Unmarshal(data, &nextbikeResp)
		if err != nil {
			return err
		}

		for _, country := range nextbikeResp.Countries {
			for _, city := range country.Cities {
				for _, place := range city.Places {
					ID := fmt.Sprintf("nextbike-%d", place.UID)

					stop := &protocol.Stop{
						ID:       ID,
						Provider: "nextbike",
						Name:     place.Name,
						Type:     "bike-stop",
						Location: protocol.Location{
							Latitude:  int(place.Lat * 3600000),
							Longitude: int(place.Lng * 3600000),
						},
					}

					d, err := json.Marshal(stop)
					if err != nil {
						return err
					}

					subject := fmt.Sprintf(protocol.SubjectMapStop, ID)
					err = c.Publish(subject, string(d))
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ Nextbike collector started")

	s.StartBlocking()
}

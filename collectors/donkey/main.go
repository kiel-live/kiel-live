package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("🚴‍♀️ Donkey collector version %s", "1.0.0") // TODO use proper version

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

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err = s.Every(5).Seconds().Do(func() error {
		if !c.IsConnected() {
			return nil
		}

		// TODO: allow to configure the bounding box
		top := "54.48855"
		left := "9.94689"
		right := "10.30319"
		bottom := "54.19533"
		url := fmt.Sprintf("https://stables.donkey.bike/api/public/nearby?top_right=%s%%2C%s&bottom_left=%s%%2C%s&filter_type=box", top, right, bottom, left)

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "donkey/1.0.0")
		req.Header.Set("Accept", "application/com.donkeyrepublic.v7")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		donkeyResp := DonkeyResponse{}
		err = json.Unmarshal(data, &donkeyResp)
		if err != nil {
			return err
		}

		for _, hub := range donkeyResp.Hubs {
			ID := fmt.Sprintf("donkey-%s", hub.ID)

			latitude, err := strconv.ParseFloat(hub.Latitude, 32)
			if err != nil {
				return err
			}

			longitude, err := strconv.ParseFloat(hub.Longitude, 32)
			if err != nil {
				return err
			}

			stop := &protocol.Stop{
				ID:       ID,
				Provider: "donkey",
				Name:     hub.Name,
				Type:     "bike-stop",
				Location: protocol.Location{
					Latitude:  int(latitude * 3600000),
					Longitude: int(longitude * 3600000),
				},
			}

			err = c.UpdateStop(stop)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ Donkey collector started")

	s.StartBlocking()
}

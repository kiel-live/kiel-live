package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

func main() {
	collector.New(collector.Options{
		Name:    "🚴‍♀️ Donkey",
		Execute: run,
	}).Run()
}

func run(_ context.Context, c client.Client) error {
	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := s.Shutdown(); err != nil {
			slog.Error("failed to shutdown scheduler", "error", err)
		}
	}()

	_, err = s.NewJob(
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(func() error {
			if !c.IsConnected() {
				return nil
			}

			// TODO: allow to configure the bounding box
			top := "54.48855"
			left := "9.94689"
			right := "10.30319"
			bottom := "54.19533"
			url := fmt.Sprintf("https://stables.donkey.bike/api/public/nearby?top_right=%s%%2C%s&bottom_left=%s%%2C%s&filter_type=box", top, right, bottom, left)

			httpClient := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Set("User-Agent", "donkey/1.0.0")
			req.Header.Set("Accept", "application/com.donkeyrepublic.v7")
			resp, err := httpClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			donkeyResp := DonkeyResponse{}
			if err = json.Unmarshal(data, &donkeyResp); err != nil {
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

				stop := &models.Stop{
					ID:       ID,
					Provider: "donkey",
					Name:     hub.Name,
					Type:     "bike-stop",
					Location: &models.Location{
						Latitude:  int(latitude * 3600000),
						Longitude: int(longitude * 3600000),
					},
				}

				if err = c.UpdateStop(stop); err != nil {
					return err
				}
			}

			return nil
		}),
	)
	if err != nil {
		return err
	}

	slog.Info("Donkey collector started")

	s.Start()
	select {}
}

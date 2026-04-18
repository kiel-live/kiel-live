package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

func main() {
	collector.New(collector.Options{
		Name:    "Nextbike",
		Execute: run,
	}).Run()
}

func run(_ context.Context, c client.Client) error {
	cityIDs := os.Getenv("NEXT_BIKE_CITY_IDS")
	if cityIDs == "" {
		return fmt.Errorf("please provide a comma separated list of next-bike city ids with NEXT_BIKE_CITY_IDS (exp: '613,195' for Kiel & Mannheim)")
	}

	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := s.Shutdown(); err != nil {
			slog.Error(err.Error())
		}
	}()

	_, err = s.NewJob(
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(func() error {
			if !c.IsConnected() {
				return nil
			}

			resp, err := http.Get("https://api.nextbike.net/maps/nextbike-live.json?city=" + cityIDs)
			if err != nil {
				return err
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			nextbikeResp := NextbikeResponse{}
			if err = json.Unmarshal(data, &nextbikeResp); err != nil {
				return err
			}

			for _, country := range nextbikeResp.Countries {
				for _, city := range country.Cities {
					for _, place := range city.Places {
						ID := fmt.Sprintf("nextbike-%d", place.UID)

						stop := &models.Stop{
							ID:       ID,
							Provider: "nextbike",
							Name:     place.Name,
							Type:     "bike-stop",
							Location: &models.Location{
								Latitude:  int(place.Lat * 3600000),
								Longitude: int(place.Lng * 3600000),
							},
							Vehicles: []*models.Vehicle{},
						}

						for _, bike := range place.BikeList {
							vehicle := &models.Vehicle{
								ID:       fmt.Sprintf("nextbike-%s", bike.Number),
								Provider: "nextbike",
								Name:     fmt.Sprintf("Nextbike %s", bike.Number),
								Type:     "bike",
								Location: &models.Location{
									Latitude:  int(place.Lat * 3600000),
									Longitude: int(place.Lng * 3600000),
								},
								State: bike.State,
								Actions: []*models.Action{
									{
										Name: "",
										Type: "rent",
										URL:  fmt.Sprintf("https://nxtb.it/%s", bike.Number),
									},
									{
										Name: "",
										Type: "navigate-to",
										URL:  fmt.Sprintf("https://www.google.com/maps/place/%f,%f", place.Lat, place.Lng),
									},
								},
								Description: "", // TODO: add pricing data
							}

							stop.Vehicles = append(stop.Vehicles, vehicle)

							if err = c.UpdateVehicle(vehicle); err != nil {
								return err
							}
						}

						if err = c.UpdateStop(stop); err != nil {
							return err
						}
					}
				}
			}

			return nil
		}),
	)
	if err != nil {
		return err
	}

	slog.Info("Nextbike collector started")

	s.Start()
	select {}
}

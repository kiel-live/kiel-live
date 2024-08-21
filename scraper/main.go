package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kiel-live/kiel-live/shared/models"
	"github.com/wI2L/jsondiff"
)

func saveDiff(time int64, lastStops []*models.Stop, stops []*models.Stop) error {
	patch, err := jsondiff.Compare(lastStops, stops)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(patch, "", "    ")
	if err != nil {
		return err
	}

	fmt.Printf("[%d] %d\n", time, len(patch))

	return os.WriteFile(fmt.Sprintf("data/diff-%d.json", time), b, 0644)
}

func main() {
	fmt.Println("Starting scraper")

	lastStops := make([]*models.Stop, 0)

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err := s.Every(5).Minutes().Do(func() error {
		stops := make([]*models.Stop, 0)

		fmt.Println("Scraping nextbike ...")
		_stops, err := nextbike()
		if err != nil {
			return err
		}
		stops = append(stops, _stops...)

		fmt.Println("Scraping donkey ...")
		_stops, err = donkey()
		if err != nil {
			return err
		}
		stops = append(stops, _stops...)

		slices.SortFunc(stops, func(a, b *models.Stop) int {
			return cmp.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))
		})

		err = saveDiff(time.Now().Unix(), lastStops, stops)
		if err != nil {
			return err
		}

		lastStops = stops
		return nil
	})
	if err != nil {
		panic(err)
	}

	s.StartBlocking()
}

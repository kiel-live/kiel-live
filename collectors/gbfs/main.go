package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
)

const IDPrefix = "gbfs-"

func main() {
	collector.New(collector.Options{
		Name:    "GBFS",
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
			slog.Error(err.Error())
		}
	}()

	_, err = s.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(func() {
			if !c.IsConnected() {
				return
			}

			// TODO load station information
			// TODO load station status
			// TODO load vehicle status
		}),
	)
	if err != nil {
		return err
	}

	slog.Info("GBFS collector started")
	s.Start()
	select {}
}

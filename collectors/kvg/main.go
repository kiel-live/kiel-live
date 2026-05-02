package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	basecollector "github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/collectors/kvg/collector"
	"github.com/kiel-live/kiel-live/pkg/client"
)

func main() {
	basecollector.New(basecollector.Options{
		Name:    "KVG",
		Execute: run,
	}).Run()
}

func run(_ context.Context, c client.Client) error {
	collectors := make(map[string]collector.Collector)

	var err error

	collectors["vehicles"], err = collector.NewCollector(c, "vehicles")
	if err != nil {
		return err
	}
	collectors["stops"], err = collector.NewCollector(c, "stops")
	if err != nil {
		return err
	}
	collectors["trips"], err = collector.NewCollector(c, "trips")
	if err != nil {
		return err
	}

	c.SetOnTopicsChanged(func(topic string, subscribed bool) {
		if !subscribed {
			return
		}

		tripID := collectors["trips"].TopicToID(topic)
		if tripID != "" {
			collectors["trips"].RunSingle(tripID)
			return
		}
		stopID := collectors["stops"].TopicToID(topic)
		if stopID != "" {
			collectors["stops"].RunSingle(stopID)
			return
		}
	})

	c.SetOnConnectionChanged(func(connected bool) {
		if !connected {
			return
		}

		for name, collector := range collectors {
			slog.Debug("Resetting collector", "name", name)
			collector.Reset()
		}
	})

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
		gocron.NewTask(func() {
			if !c.IsConnected() {
				return
			}

			for name, collector := range collectors {
				// TODO maybe run in go routine
				slog.Debug("Running collector ...", "name", name)
				collector.Run()
			}
		}),
	)
	if err != nil {
		return err
	}

	slog.Info("KVG collector started")

	s.Start()
	select {}
}

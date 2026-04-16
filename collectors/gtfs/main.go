package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/artonge/go-gtfs"
	"github.com/go-co-op/gocron/v2"
	"github.com/hashicorp/go-memdb"

	"github.com/kiel-live/kiel-live/collectors/collector"
	"github.com/kiel-live/kiel-live/collectors/gtfs/loader"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
)

const IDPrefix = "gtfs-"

func main() {
	collector.New(collector.Options{
		Name:    "GTFS",
		Execute: run,
	}).Run()
}

func run(_ context.Context, c client.Client) error {
	// Example: https://github.com/lukashass/nok-gtfs/raw/main/adler.zip
	gtfsPath := os.Getenv("GTFS_PATH")
	if gtfsPath == "" {
		return fmt.Errorf("please provide a GTFS path with GTFS_PATH")
	}

	// ; separated list of alerts
	generalAlerts := os.Getenv("GTFS_GENERAL_ALERTS")

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"stop_times": {
				Name: "stop_times",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "TripID"},
								&memdb.StringFieldIndex{Field: "StopID"},
								&memdb.UintFieldIndex{Field: "StopSeq"},
							},
						},
					},
					"stop_id": {
						Name:    "stop_id",
						Indexer: &memdb.StringFieldIndex{Field: "StopID"},
					},
					"trip_id": {
						Name:    "trip_id",
						Indexer: &memdb.StringFieldIndex{Field: "TripID"},
					},
				},
			},
			"trips": {
				Name: "trips",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"route_id": {
						Name:    "route_id",
						Indexer: &memdb.StringFieldIndex{Field: "RouteID"},
					},
					"service_id": {
						Name:    "service_id",
						Indexer: &memdb.StringFieldIndex{Field: "ServiceID"},
					},
				},
			},
			"routes": {
				Name: "routes",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
			"calendars": {
				Name: "calendars",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ServiceID"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return fmt.Errorf("failed opening db: %w", err)
	}

	g, err := loader.LoadGTFS(gtfsPath)
	if err != nil {
		return err
	}

	agency := g.Agency
	if len(g.Agencies) > 1 {
		return fmt.Errorf("multiple agencies are not supported")
	}

	// create stop times table
	txn := db.Txn(true)

	for _, stopTime := range g.StopsTimes {
		if err = txn.Insert("stop_times", stopTime); err != nil {
			return err
		}
	}

	for _, trip := range g.Trips {
		if err = txn.Insert("trips", trip); err != nil {
			return err
		}

		// delete the last stop time for each trip (highest stop_sequence),
		// because the vehicle does not depart in this trip anymore, the trip is finished
		stopTimesIt, err := txn.Get("stop_times", "trip_id", trip.ID)
		if err != nil {
			return err
		}

		var lastStopTime *gtfs.StopTime
		for obj := stopTimesIt.Next(); obj != nil; obj = stopTimesIt.Next() {
			stopTime := obj.(gtfs.StopTime)

			if lastStopTime == nil {
				lastStopTime = &stopTime
			} else if stopTime.StopSeq > lastStopTime.StopSeq {
				lastStopTime = &stopTime
			}
		}

		if lastStopTime != nil {
			if err = txn.Delete("stop_times", lastStopTime); err != nil {
				return err
			}
		}
	}

	for _, route := range g.Routes {
		if err = txn.Insert("routes", route); err != nil {
			return err
		}
	}

	for _, calendar := range g.Calendars {
		if err = txn.Insert("calendars", calendar); err != nil {
			return err
		}
	}

	txn.Commit()

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

			txn := db.Txn(false)
			defer txn.Abort()

			for _, gtfsStop := range g.Stops {
				// convert to protocol.Stop
				stop := &models.Stop{
					ID:       IDPrefix + gtfsStop.ID,
					Provider: agency.Name,
					Name:     gtfsStop.Name,
					Type:     "",
					Location: &models.Location{
						Latitude:  int(gtfsStop.Latitude * 3600000),
						Longitude: int(gtfsStop.Longitude * 3600000),
					},
					Alerts: strings.Split(generalAlerts, ";"), // TODO: get alerts from gtfs-rt feed
				}

				// TODO: remove empty alerts
				if len(stop.Alerts) == 1 && stop.Alerts[0] == "" {
					stop.Alerts = []string{}
				}

				stopTimesIt, err := txn.Get("stop_times", "stop_id", gtfsStop.ID)
				if err != nil {
					slog.Error(err.Error())
					continue
				}

				stopTimes := make([]gtfs.StopTime, 0)
				for obj := stopTimesIt.Next(); obj != nil; obj = stopTimesIt.Next() {
					stopTime := obj.(gtfs.StopTime)
					stopTimes = append(stopTimes, stopTime)
				}

				// iterate over stop times
				for _, stopTime := range stopTimes {
					_trip, err := txn.First("trips", "id", stopTime.TripID)
					if err != nil {
						slog.Error(err.Error())
						continue
					}
					trip := _trip.(gtfs.Trip)

					_calendar, err := txn.First("calendars", "id", trip.ServiceID)
					if err != nil {
						slog.Error(err.Error())
						continue
					}
					calendar := _calendar.(gtfs.Calendar)

					// check if service is active today
					if !weekdayIsActiveInCalendar(calendar) {
						continue
					}

					_route, err := txn.First("routes", "id", trip.RouteID)
					if err != nil {
						slog.Error(err.Error())
						continue
					}
					route := _route.(gtfs.Route)

					// convert departure time to unix timestamp
					departureTime, err := time.Parse("15:04:05", stopTime.Departure)
					if err != nil {
						slog.Error(err.Error())
						continue
					}

					now := time.Now()
					departureDate := time.Date(now.Year(), now.Month(), now.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), 0, time.Local)

					if departureDate.Before(now) || departureDate.After(now.Add(4*time.Hour)) {
						continue
					}

					// TODO: consider all routes for stop type
					if stop.Type == "" {
						stop.Type = models.StopType(gtfsRouteTypeToProtocolStopType(route.Type) + "-stop")
					}

					stop.Departures = append(stop.Departures, &models.StopDepartures{
						Name:      stop.Name,
						Type:      models.VehicleType(gtfsRouteTypeToProtocolStopType(route.Type)),
						TripID:    IDPrefix + stopTime.TripID,
						Actual:    "", // TODO: get from gtfs-rt
						Planned:   departureDate.Format(time.RFC3339),
						RouteName: route.ShortName,
						Direction: trip.Headsign,
						State:     models.Planned,
						RouteID:   IDPrefix + route.ID,
						VehicleID: IDPrefix + trip.ID,
						Platform:  "", // TODO
					})
				}
				if stop.Type == "" {
					slog.Warn("Stop has no type and is therefore skipped", "stop_id", stop.ID)
					continue
				}

				if len(stop.Departures) == 0 {
					slog.Warn("Stop has no departures and is therefore skipped", "stop_id", stop.ID)
					continue
				}

				if err = c.UpdateStop(stop); err != nil {
					slog.Error(err.Error())
					continue
				}
			}
		}),
	)
	if err != nil {
		return err
	}

	slog.Info("GTFS collector started")

	s.Start()
	select {}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/artonge/go-gtfs"
	"github.com/go-co-op/gocron/v2"
	"github.com/hashicorp/go-memdb"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/gtfs/loader"
	"github.com/kiel-live/kiel-live/collectors/gtfs/rt"
	"github.com/kiel-live/kiel-live/protocol"
)

const IDPrefix = "gtfs-"

func main() {
	log.Infof("Kiel-Live GTFS collector version %s", "1.0.0") // TODO use proper version

	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Fatalf("error loading location '%s': %v\n", tz, err)
		}
	}

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

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		log.Fatalln("Please provide a token for the collector with MANAGER_TOKEN")
	}

	// Example: https://github.com/lukashass/nok-gtfs/raw/main/adler.zip
	gtfsPath := os.Getenv("GTFS_PATH")
	if gtfsPath == "" {
		log.Fatalln("Please provide a GTFS path with GTFS_PATH")
	}

	gtfsRT, err := rt.NewGTFSRTCollector(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	// ; separated list of alerts
	generalAlerts := os.Getenv("GTFS_GENERAL_ALERTS")

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

	g, err := loader.LoadGTFS(gtfsPath)
	if err != nil {
		log.Error(err)
		return
	}

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
		log.Panic("Failed opening db", err)
	}

	agency := g.Agency
	if g.Agencies != nil && len(g.Agencies) > 1 {
		log.Fatal("Multiple agencies are not supported")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	// update gtfs feed
	_, err = s.NewJob(gocron.DurationJob(time.Hour*24), gocron.NewTask(func() {
		g, err := loader.LoadGTFS(gtfsPath)
		if err != nil {
			log.Error(err)
			return
		}

		err = importGTFS(db, g)
		if err != nil {
			log.Error(err)
			return
		}
	}), gocron.WithLimitedRuns(1), gocron.JobOption(gocron.WithStartImmediately()))
	if err != nil {
		log.Fatal(err)
	}

	// update gtfs-rt feed
	_, err = s.NewJob(gocron.DurationJob(time.Minute*1), gocron.NewTask(func() {
		if gtfsRT == nil {
			return
		}

		txn := db.Txn(false)
		defer txn.Abort()

		feed, err := gtfsRT.FetchTripUpdates()
		if err != nil {
			log.Error(err)
		}

		if feed == nil {
			return
		}

		// TODO: nil check
		for _, entity := range feed.Entity {
			tripUpdate := entity.GetTripUpdate()

			stopTimesIt, err := txn.Get("stop_times", "trip_id", tripUpdate.Trip.TripId)
			if err != nil {
				log.Error(err)
				return
			}

			stopTimes := make(map[uint32]*gtfs.StopTime)
			for obj := stopTimesIt.Next(); obj != nil; obj = stopTimesIt.Next() {
				stopTime := obj.(gtfs.StopTime)
				stopTimes[stopTime.StopSeq] = &stopTime
			}

			// todo add to db or update somehow
			for _, u := range tripUpdate.StopTimeUpdate {
				stopTime, ok := stopTimes[u.GetStopSequence()]
				if !ok {
					log.Warnf("Stop time %d not found for trip %s", u.GetStopSequence(), *tripUpdate.Trip.TripId)
					continue
				}

				// update stop time
				// stopTime.Arrival = time.Unix(int64(u.Arrival.Time), 0).Format("15:04:05")
				// stopTime.Departure = u.Departure.GetTime()

				err := txn.Insert("stop_times", stopTime)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}), gocron.WithLimitedRuns(1), gocron.JobOption(gocron.WithStartImmediately()))
	if err != nil {
		log.Fatal(err)
	}

	// publish stops and arrivals
	_, err = s.NewJob(gocron.DurationJob(time.Minute*1), gocron.NewTask(func() {
		if !c.IsConnected() {
			return
		}

		txn := db.Txn(false)
		defer txn.Abort()

		stopIt, err := txn.Get("stops", "")
		if err != nil {
			log.Error(err)
			return
		}

		for obj := stopIt.Next(); obj != nil; obj = stopIt.Next() {
			gtfsStop := obj.(gtfs.Stop)

			// convert to protocol.Stop
			stop := protocol.Stop{
				ID:       IDPrefix + gtfsStop.ID,
				Provider: agency.Name,
				Name:     gtfsStop.Name,
				Type:     "",
				Location: protocol.Location{
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
				log.Error(err)
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
					log.Error(err)
					continue
				}
				trip := _trip.(gtfs.Trip)

				_calendar, err := txn.First("calendars", "id", trip.ServiceID)
				if err != nil {
					log.Error(err)
					continue
				}
				calendar := _calendar.(gtfs.Calendar)

				// check if service is active today
				if !weekdayIsActiveInCalendar(calendar) {
					continue
				}

				_route, err := txn.First("routes", "id", trip.RouteID)
				if err != nil {
					log.Error(err)
					continue
				}
				route := _route.(gtfs.Route)

				// convert departure time to unix timestamp
				departureTime, err := time.Parse("15:04:05", stopTime.Departure)
				if err != nil {
					log.Error(err)
					continue
				}

				now := time.Now()
				departureDate := time.Date(now.Year(), now.Month(), now.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), 0, time.Local)
				if departureDate.Before(now) || departureDate.After(now.Add(4*time.Hour)) {
					continue
				}

				// TODO: consider all routes for stop type
				if stop.Type == "" {
					stop.Type = protocol.StopType(gtfsRouteTypeToProtocolStopType(route.Type) + "-stop")
				}

				stop.Arrivals = append(stop.Arrivals, protocol.StopArrival{
					Name:      stop.Name,
					Type:      protocol.VehicleType(gtfsRouteTypeToProtocolStopType(route.Type)),
					TripID:    IDPrefix + stopTime.TripID,
					ETA:       0, // TODO: get from gtfs-rt
					Planned:   departureDate.Format("15:04"),
					RouteName: route.ShortName,
					Direction: trip.Headsign,
					State:     protocol.Planned,
					RouteID:   IDPrefix + route.ID,
					VehicleID: IDPrefix + trip.ID,
					Platform:  "", // TODO
				})
			}
			if stop.Type == "" {
				log.Warnf("Stop %s has no type and is therefore skipped", stop.ID)
				continue
			}

			if len(stop.Arrivals) == 0 {
				log.Warnf("Stop %s has no arrivals and is therefore skipped", stop.ID)
				continue
			}

			jsonData, err := json.Marshal(stop)
			if err != nil {
				log.Error(err)
				continue
			}

			// publish stop
			subject := fmt.Sprintf(protocol.SubjectMapStop, stop.ID)
			err = c.Publish(subject, string(jsonData))
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}), gocron.WithLimitedRuns(1), gocron.JobOption(gocron.WithStartImmediately()))
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("⚡ GTFS collector started")

	s.Start()

	// block forever
	select {}
}

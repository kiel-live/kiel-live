package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/artonge/go-gtfs"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/collectors/gtfs/loader"
	"github.com/kiel-live/kiel-live/protocol"
	log "github.com/sirupsen/logrus"
)

const IDPrefix = "gtfs-"

func main() {
	log.Infof("Kiel-Live GTFS collector version %s", "1.0.0") // TODO use proper version

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

	generalAlerts := os.Getenv("GENERAL_ALERTS")

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

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err = s.Every(1).Minute().Do(func() {
		if !c.IsConnected() {
			return
		}

		agency := g.Agency
		if g.Agencies != nil && len(g.Agencies) > 1 {
			log.Fatal("Multiple agencies are not supported")
		}

		for _, gtfsStop := range g.Stops {
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

			// for each trip remove stop times with the highest stop_sequence (last stop)
			// TODO: solve differently to avoid additional O(n^2)
			// filteredStopTimes := make([]gtfs.StopTime, 0)
			// for _, stopTime := range g.StopsTimes {
			// 	isLastStop := true
			// 	for _, stopTime2 := range g.StopsTimes {
			// 		if stopTime.TripID == stopTime2.TripID && stopTime.StopSeq < stopTime2.StopSeq {
			// 			isLastStop = false
			// 			break
			// 		}
			// 	}
			// 	if !isLastStop {
			// 		filteredStopTimes = append(filteredStopTimes, stopTime)
			// 	}
			// }

			// iterate over stop times
			for _, stopTime := range g.StopsTimes {
				if stopTime.StopID == gtfsStop.ID {
					index, found := findInObjArr(g.Trips, func(t gtfs.Trip) string { return t.ID }, stopTime.TripID)
					if !found {
						log.Warnf("Trip %s not found", stopTime.TripID)
						continue
					}
					trip := g.Trips[index]

					index, found = findInObjArr(g.Calendars, func(c gtfs.Calendar) string { return c.ServiceID }, trip.ServiceID)
					if !found {
						log.Warnf("Calendar %s not found", trip.ServiceID)
						continue
					}
					calendar := g.Calendars[index]
					// check if service is active today
					if !weekdayIsActiveInCalendar(calendar) {
						continue
					}

					index, found = findInObjArr(g.Routes, func(r gtfs.Route) string { return r.ID }, trip.RouteID)
					if !found {
						log.Warnf("Route %s not found", trip.RouteID)
						continue
					}
					route := g.Routes[index]

					// TODO: consider all routes for stop type
					stop.Type = protocol.StopType(gtfsRouteTypeToProtocolStopType(route.Type) + "-stop")

					// convert departure time to unix timestamp
					departureTime, err := time.Parse("15:04:05", stopTime.Departure)
					if err != nil {
						// log.Debug(err)
						continue
					}
					now := time.Now()
					departureDate := time.Date(now.Year(), now.Month(), now.Day(), departureTime.Hour(), departureTime.Minute(), departureTime.Second(), 0, time.Local)
					if departureDate.Before(now) || departureDate.After(now.Add(4*time.Hour)) {
						continue
					}

					stop.Arrivals = append(stop.Arrivals, protocol.StopArrival{
						Name:      stop.Name,
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
			}
			if stop.Type == "" {
				// log.Warnf("Stop %s has no type and is therefore skipped", stop.ID)
				continue
			}

			if len(stop.Arrivals) == 0 {
				// log.Warnf("Stop %s has no arrivals and is therefore skipped", stop.ID)
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
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("⚡ GTFS collector started")

	s.StartBlocking()
}

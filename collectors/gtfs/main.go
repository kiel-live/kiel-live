package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/artonge/go-gtfs"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/client"
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

	// load gtfs files
	g, err := gtfs.Load("collectors/gtfs/adler", nil)
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

		for _, gtfsStop := range g.Stops {
			// convert to protocol.Stop
			stop := protocol.Stop{
				ID:   IDPrefix + gtfsStop.ID,
				Name: gtfsStop.Name,
				Type: protocol.StopTypeFerryStop,
				Location: protocol.Location{
					Latitude:  int(gtfsStop.Latitude * 3600000),
					Longitude: int(gtfsStop.Longitude * 3600000),
				},
				Alerts: []string{"Die Abfahrtszeiten können sich je nach Witterung oder Verkehrslage auf dem Nord-Ostsee-Kanal geringfügig verschieben. Die Verschiebung einer Fahrt dient der Sicherheit des Fahrbetriebes. Bei Ausfall der Fähre ist ein Busersatzverkehr eingerichtet."},
			}
			// for each trip remove stop times with the highest stop_sequence (last stop)
			filteredStopTimes := make([]gtfs.StopTime, 0)
			for _, stopTime := range g.StopsTimes {
				isLastStop := true
				for _, stopTime2 := range g.StopsTimes {
					if stopTime.TripID == stopTime2.TripID && stopTime.StopSeq < stopTime2.StopSeq {
						isLastStop = false
						break
					}
				}
				if !isLastStop {
					filteredStopTimes = append(filteredStopTimes, stopTime)
				}
			}

			// iterate over stop times
			for _, stopTime := range filteredStopTimes {
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

					index, found = findInObjArr(g.Routes, func(r gtfs.Route) string { return r.ID }, trip.RouteID)
					if !found {
						log.Warnf("Route %s not found", trip.RouteID)
						continue
					}
					route := g.Routes[index]

					// check if service is active
					// get current weekday
					if !weekdayIsActiveInCalendar(calendar) {
						continue
					}

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

					stop.Arrivals = append(stop.Arrivals, protocol.StopArrival{
						TripID:    IDPrefix + stopTime.TripID,
						Planned:   departureDate.Format("15:04"),
						RouteName: route.ShortName,
						Direction: trip.Headsign,
						State:     protocol.Planned,
					})
				}
			}
			jsonData, err := json.Marshal(stop)
			if err != nil {
				log.Error(err)
				continue
			}

			subject := fmt.Sprintf(protocol.SubjectMapStop, stop.ID)

			// publish stop
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

func findInObjArr[T any, K comparable](arr []T, keyFunc func(T) K, value K) (int, bool) {
	for i, v := range arr {
		if keyFunc(v) == value {
			return i, true
		}
	}
	return -1, false
}

func weekdayIsActiveInCalendar(calendar gtfs.Calendar) bool {
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Monday:
		return calendar.Monday == 1
	case time.Tuesday:
		return calendar.Tuesday == 1
	case time.Wednesday:
		return calendar.Wednesday == 1
	case time.Thursday:
		return calendar.Thursday == 1
	case time.Friday:
		return calendar.Friday == 1
	case time.Saturday:
		return calendar.Saturday == 1
	case time.Sunday:
		return calendar.Sunday == 1
	}
	return false
}

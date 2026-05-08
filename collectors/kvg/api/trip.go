package api

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/kiel-live/kiel-live/pkg/models"
)

type tripItem struct {
	Stop         tripStop        `json:"stop"`
	Status       DepartureStatus `json:"status"`
	ActualTime   string          `json:"actualTime"`
	StopSequence string          `json:"stop_seq_num"`
}

type tripStop struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type waypoint struct {
	Lat int `json:"lat"`
	Lon int `json:"lon"`
}

type tripPath struct {
	Waypoints []waypoint `json:"wayPoints"`
}

type tripPaths struct {
	Paths []tripPath `json:"paths"`
}

func (t *tripItem) parse() *models.TripDeparture {
	actual, err := timeToIsoDateTime(t.ActualTime, time.Now())
	if err != nil {
		actual = ""
		log.Printf("Error parsing time for trip departure: %v", err)
	}

	return &models.TripDeparture{
		ID:      IDPrefix + t.Stop.ShortName,
		Name:    t.Stop.Name,
		State:   t.Status.parse(),
		Planned: actual, // KVG API does not provide planned time, using actual time as fallback
		Actual:  actual,
	}
}

type trip struct {
	Stops         []tripItem `json:"actual"`
	OldStops      []tripItem `json:"old"`
	DirectionText string     `json:"directionText"`
	RouteName     string     `json:"routeName"`
}

// trip parser to protocol trip
func (t *trip) parse() *models.Trip {
	var departures []*models.TripDeparture
	for _, stop := range t.OldStops {
		departures = append(departures, stop.parse())
	}
	for _, stop := range t.Stops {
		departures = append(departures, stop.parse())
	}

	return &models.Trip{
		Provider:   "kvg", // TODO
		Direction:  t.DirectionText,
		Departures: departures,
	}
}

func GetTripPath(tripID string) []*models.Location {
	data := url.Values{}
	data.Set("id", tripID)

	resp, err := post(tripPathURL, data)
	if err != nil {
		return nil
	}
	var paths tripPaths
	if err := json.Unmarshal(resp, &paths); err != nil {
		return nil
	}

	var path []*models.Location
	for _, waypoint := range paths.Paths[0].Waypoints {
		path = append(path, &models.Location{
			Latitude:  waypoint.Lat,
			Longitude: waypoint.Lon,
		})
	}

	return path
}

func GetTrip(tripID string) (*models.Trip, error) {
	data := url.Values{}
	data.Set("tripId", tripID)

	resp, err := post(tripURL, data)
	if err != nil {
		return nil, err
	}

	var trip trip
	if err := json.Unmarshal(resp, &trip); err != nil {
		return nil, err
	}

	protocolTrip := trip.parse()
	protocolTrip.ID = IDPrefix + tripID
	protocolTrip.Path = GetTripPath(tripID)

	return protocolTrip, nil
}

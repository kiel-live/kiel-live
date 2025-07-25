package api

import (
	"encoding/json"
	"net/url"

	"github.com/kiel-live/kiel-live/pkg/models"
)

type tripStop struct {
	Stop       stop            `json:"stop"`
	Status     DepartureStatus `json:"status"`
	ActualTime string          `json:"actualTime"`
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

func (t *tripStop) parse() *models.TripArrival {
	return &models.TripArrival{
		ID:      IDPrefix + t.Stop.ShortName,
		Name:    t.Stop.Name,
		State:   t.Status.parse(),
		Planned: t.ActualTime,
	}
}

type trip struct {
	Stops         []tripStop `json:"actual"`
	OldStops      []tripStop `json:"old"`
	DirectionText string     `json:"directionText"`
	RouteName     string     `json:"routeName"`
}

// trip parser to protocol trip
func (t *trip) parse() *models.Trip {
	var arrivals []*models.TripArrival
	for _, stop := range t.OldStops {
		arrivals = append(arrivals, stop.parse())
	}
	for _, stop := range t.Stops {
		arrivals = append(arrivals, stop.parse())
	}

	return &models.Trip{
		Provider:  "kvg", // TODO
		Direction: t.DirectionText,
		Arrivals:  arrivals,
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

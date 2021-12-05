package api

import (
	"encoding/json"
	"net/url"

	"github.com/kiel-live/kiel-live/protocol"
)

type tripStop struct {
	Stop       stop            `json:"stop"`
	Status     DepartureStatus `json:"status"`
	ActualTime string          `json:"actualTime"`
}

func (t *tripStop) parse() protocol.TripArrival {
	return protocol.TripArrival{
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
func (t *trip) parse() protocol.Trip {
	var arrivals []protocol.TripArrival
	for _, stop := range t.OldStops {
		arrivals = append(arrivals, stop.parse())
	}
	for _, stop := range t.Stops {
		arrivals = append(arrivals, stop.parse())
	}

	return protocol.Trip{
		Provider:  "kvg", // TODO
		Direction: t.DirectionText,
		Arrivals:  arrivals,
	}
}

func GetTrip(tripID string) (protocol.Trip, error) {
	data := url.Values{}
	data.Set("tripId", tripID)

	resp, _ := post(tripURL, data)
	var trip trip
	if err := json.Unmarshal(resp, &trip); err != nil {
		return protocol.Trip{}, err
	}
	protocolTrip := trip.parse()
	protocolTrip.ID = IDPrefix + tripID
	return protocolTrip, nil
}

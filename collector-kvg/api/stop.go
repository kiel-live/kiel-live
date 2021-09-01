package api

import (
	"encoding/json"
	"log"
	"net/url"

	protocol "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
)

type stop struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
	Latitude  int    `json:"latitude"`
	Longitude int    `json:"longitude"`
}

func (s *stop) parse() protocol.Stop {
	return protocol.Stop{
		ID:       s.ID,
		Provider: "kvg", // TODO
		Name:     s.Name,
		Type:     protocol.VehicleTypeBus,
		Location: protocol.Location{
			Longitude: float32(s.Longitude),
			Latitude:  float32(s.Latitude),
		},
	}
}

type stops struct {
	Stops []stop `json:"stops"`
}

type departure struct {
	TripID             string `json:"tripId"`
	Status             string `json:"status"`
	Stop               string `json:"plannedTime"`
	ActualTime         string `json:"actualTime"`
	ActualRelativeTime int    `json:"actualRelativeTime"`
}

type stopDepartures struct {
	Departures []departure `json:"actual"`
}

func GetStops() []protocol.Stop {
	data := url.Values{}
	data.Set("top", "324000000")
	data.Set("bottom", "-324000000")
	data.Set("left", "-648000000")
	data.Set("right", "648000000")

	body, _ := post(stopsURL, data)
	var rawStops stops
	if err := json.Unmarshal(body, &rawStops); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}

	var stops []protocol.Stop
	for _, stop := range rawStops.Stops {
		stops = append(stops, stop.parse())
	}

	return stops
}

func GetStop(stopShortName string) stopDepartures {
	data := url.Values{}
	data.Set("stop", stopShortName)

	resp, _ := post(stopURL, data)
	var stop stopDepartures
	if err := json.Unmarshal(resp, &stop); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	return stop
}

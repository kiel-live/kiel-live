package api

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/kiel-live/kiel-live/protocol"
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
		ID:       "kvg-" + s.ShortName,
		Provider: "kvg", // TODO
		Name:     s.Name,
		Type:     protocol.VehicleTypeBus,
		Location: protocol.Location{
			Longitude: s.Longitude,
			Latitude:  s.Latitude,
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
	VehicleID          string `json:"vehicleId"`
	RouteID            string `json:"routeId"`
	Direction          string `json:"direction"`
}

func (d *departure) parse() protocol.StopArrival {
	return protocol.StopArrival{
		Name:      d.Direction,
		VehicleID: d.VehicleID,
		TripID:    d.TripID,
		RouteID:   d.RouteID,
		Direction: d.Direction,
		State:     d.Status,
		ETA:       d.ActualRelativeTime,
	}
}

type StopDepartures struct {
	Departures []departure `json:"actual"`
}

func GetStops() (res map[string]*protocol.Stop) {
	res = make(map[string]*protocol.Stop)
	data := url.Values{}
	data.Set("top", "324000000")
	data.Set("bottom", "-324000000")
	data.Set("left", "-648000000")
	data.Set("right", "648000000")

	body, _ := post(stopsURL, data)
	var stops stops
	if err := json.Unmarshal(body, &stops); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}

	for _, stop := range stops.Stops {
		v := stop.parse()
		res[v.ID] = &v
	}

	return res
}

func GetStopDepartures(stopShortName string) []protocol.StopArrival {
	data := url.Values{}
	data.Set("stop", stopShortName)

	resp, _ := post(stopURL, data)
	var stop StopDepartures
	if err := json.Unmarshal(resp, &stop); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n response: %v", err, string(resp))
	}

	departures := []protocol.StopArrival{}
	for _, departure := range stop.Departures {
		departures = append(departures, departure.parse())
	}
	return departures
}

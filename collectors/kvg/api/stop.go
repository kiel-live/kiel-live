package api

import (
	"encoding/json"
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
		ID:       IDPrefix + s.ShortName,
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

type DepartureStatus string

const (
	planned   DepartureStatus = "PLANNED"
	predicted                 = "PREDICTED"
	stopping                  = "STOPPING"
)

func (d *DepartureStatus) parse() protocol.ArrivalState {
	switch *d {
	case planned:
		return protocol.Planned
	case predicted:
		return protocol.Predicted
	case stopping:
		return protocol.Stopping
	default:
		return protocol.Undefined
	}
}

type departure struct {
	TripID             string          `json:"tripId"`
	Status             DepartureStatus `json:"status"`
	Stop               string          `json:"plannedTime"`
	ActualTime         string          `json:"actualTime"`
	ActualRelativeTime int             `json:"actualRelativeTime"`
	VehicleID          string          `json:"vehicleId"`
	RouteID            string          `json:"routeId"`
	RouteName          string          `json:"patternText"`
	Direction          string          `json:"direction"`
}

func (d *departure) parse() protocol.StopArrival {
	return protocol.StopArrival{
		Name:      d.Direction,
		VehicleID: d.VehicleID,
		TripID:    d.TripID,
		RouteID:   d.RouteID,
		RouteName: d.RouteName,
		Direction: d.Direction,
		State:     d.Status.parse(),
		ETA:       d.ActualRelativeTime,
		Planned:   d.ActualTime,
	}
}

type StopDepartures struct {
	Departures []departure `json:"actual"`
}

func GetStops() (res map[string]*protocol.Stop, err error) {
	res = make(map[string]*protocol.Stop)
	data := url.Values{}
	data.Set("top", "324000000")
	data.Set("bottom", "-324000000")
	data.Set("left", "-648000000")
	data.Set("right", "648000000")

	body, _ := post(stopsURL, data)
	var stops stops
	if err := json.Unmarshal(body, &stops); err != nil {
		return nil, err
	}

	for _, stop := range stops.Stops {
		v := stop.parse()
		res[v.ID] = &v
	}

	return res, nil
}

func GetStopDepartures(stopShortName string) ([]protocol.StopArrival, error) {
	data := url.Values{}
	data.Set("stop", stopShortName)

	resp, _ := post(stopURL, data)
	var stop StopDepartures
	if err := json.Unmarshal(resp, &stop); err != nil {
		return nil, err
	}

	departures := []protocol.StopArrival{}
	for _, departure := range stop.Departures {
		departures = append(departures, departure.parse())
	}
	return departures, nil
}

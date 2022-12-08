package api

import (
	"encoding/json"
	"net/url"

	"github.com/kiel-live/kiel-live/protocol"
)

type stop struct {
	ID        string   `json:"id"`
	ShortName string   `json:"shortName"`
	Name      string   `json:"name"`
	Latitude  int      `json:"latitude"`
	Longitude int      `json:"longitude"`
	Alerts    []string `json:"alerts"`
}

func (s *stop) parse() protocol.Stop {
	return protocol.Stop{
		ID:       IDPrefix + s.ShortName,
		Provider: "kvg", // TODO
		Name:     s.Name,
		Type:     protocol.StopTypeBusStop,
		Alerts:   s.Alerts,
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
	predicted DepartureStatus = "PREDICTED"
	stopping  DepartureStatus = "STOPPING"
	departed  DepartureStatus = "DEPARTED"
)

func (d *DepartureStatus) parse() protocol.ArrivalState {
	switch *d {
	case planned:
		return protocol.Planned
	case predicted:
		return protocol.Predicted
	case stopping:
		return protocol.Stopping
	case departed:
		return protocol.Departed
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

type alert struct {
	Title string `json:"title"`
}

type routes struct {
	Alerts     []alert  `json:"alerts"`
	Authority  string   `json:"authority"`
	Directions []string `json:"directions"`
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	RouteType  string   `json:"routeType"`
	ShortName  string   `json:"shortName"`
}

func (d *departure) parse() protocol.StopArrival {
	return protocol.StopArrival{
		Name:      d.Direction,
		VehicleID: IDPrefix + d.VehicleID,
		TripID:    IDPrefix + d.TripID,
		RouteID:   d.RouteID,
		RouteName: d.RouteName,
		Direction: d.Direction,
		State:     d.Status.parse(),
		ETA:       d.ActualRelativeTime,
		Planned:   d.ActualTime,
	}
}

type StopDepartures struct {
	Departures    []departure `json:"actual"`
	GeneralAlerts []alert     `json:"generalAlerts"`
	Routes        []routes    `json:"routes"`
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

type StopDetails struct {
	Departures []protocol.StopArrival
	Alerts     []string
}

func GetStopDetails(stopShortName string) (*StopDetails, error) {
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

	alerts := []string{}
	for _, alert := range stop.GeneralAlerts {
		alerts = append(alerts, alert.Title)
	}
	for _, route := range stop.Routes {
		for _, alert := range route.Alerts {
			alerts = append(alerts, route.Name+": "+alert.Title)
		}
	}

	alerts = removeDuplicate(alerts)

	details := &StopDetails{
		Departures: departures,
		Alerts:     alerts,
	}

	return details, nil
}

// https://stackoverflow.com/a/66751055/14671646
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

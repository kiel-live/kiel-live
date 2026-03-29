package models

import "fmt"

// Type of a stop
type StopType string

const (
	StopTypeBusStop     StopType = "bus-stop"
	StopTypeParkingSpot StopType = "parking-spot"
	StopTypeFerryStop   StopType = "ferry-stop"
	StopTypeTrainStop   StopType = "train-stop"
	StopTypeSubwayStop  StopType = "subway-stop"
)

// Stop is a fixed point where for example a bus stops or a car-sharing parking spot is located
type Stop struct {
	ID         string            `json:"id"`
	Provider   string            `json:"provider"`
	Name       string            `json:"name"`
	Type       StopType          `json:"type"`             // Deprecated: use departures[].type or vehicles[].type instead
	Routes     []*Route          `json:"routes,omitempty"` // list of routes using this stop
	Alerts     []string          `json:"alerts,omitempty"` // general alerts for this stop
	Departures []*StopDepartures `json:"departures"`       // omitempty is avoided as null / [] is used to check if data was loaded
	Vehicles   []*Vehicle        `json:"vehicles,omitempty"`
	Location   *Location         `json:"location"`
}

func (s *Stop) String() string {
	if s == nil {
		return "Stop(nil)"
	}
	return fmt.Sprintf("Stop(%s, %s, %s)", s.ID, s.Provider, s.Name)
}

type DepartureState string

// State of a departure
const (
	Undefined DepartureState = ""
	Stopping  DepartureState = "stopping"
	Predicted DepartureState = "predicted"
	Planned   DepartureState = "planned"
	Departed  DepartureState = "departed"
)

type StopDepartures struct {
	Name      string         `json:"name"`
	Type      VehicleType    `json:"type"`
	VehicleID string         `json:"vehicleId"`
	TripID    string         `json:"tripId"`
	RouteID   string         `json:"routeId"`
	RouteName string         `json:"routeName"`
	Direction string         `json:"direction"`
	State     DepartureState `json:"state"`
	Planned   string         `json:"planned"`
	Actual    string         `json:"actual,omitempty"`
	// Eta       int            `json:"eta"`
	Platform string `json:"platform"`
}

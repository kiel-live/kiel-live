package protocol

import (
	"time"
)

type StopArrivalState string

// State of an arrival
const (
	Undefined StopArrivalState = ""
	Stopping                   = "stopping"
	Predicted                  = "predicted"
	Planned                    = "planned"
)

type StopArrival struct {
	Name      string           `json:"name"`
	VehicleID string           `json:"vehicleId"`
	TripID    string           `json:"tripId"`
	RouteID   string           `json:"routeId"`
	RouteName string           `json:"routeName"`
	Direction string           `json:"direction"`
	State     StopArrivalState `json:"state"`
	Planned   time.Time        `json:"planned"`
	ETA       int              `json:"eta"` // in seconds
}

type TripArrival struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	State   string    `json:"state"` // use ArrivalState...
	Planned time.Time `json:"planned"`
}

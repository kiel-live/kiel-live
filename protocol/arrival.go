package protocol

import "time"

// State of an arrival
const (
	ArrivalStatePredicted = "predicted"
	ArrivalStateStopping  = "stopping"
	ArrivalStateDeparted  = "departed"
)

type StopArrival struct {
	Name      string    `json:"name"`
	VehicleID string    `json:"vehicleId"`
	TripID    string    `json:"tripId"`
	RouteID   string    `json:"routeId"`
	Direction string    `json:"direction"`
	State     string    `json:"state"` // use ArrivalState...
	Planned   time.Time `json:"planned"`
	ETA       int       `json:"eta"` // in seconds
}

type TripArrival struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	State   string    `json:"state"` // use ArrivalState...
	Planned time.Time `json:"planned"`
}

package protocol

type ArrivalState string

// State of an arrival
const (
	Undefined ArrivalState = ""
	Stopping               = "stopping"
	Predicted              = "predicted"
	Planned                = "planned"
)

type StopArrival struct {
	Name      string       `json:"name"`
	VehicleID string       `json:"vehicleId"`
	TripID    string       `json:"tripId"`
	RouteID   string       `json:"routeId"`
	RouteName string       `json:"routeName"`
	Direction string       `json:"direction"`
	State     ArrivalState `json:"state"`
	Planned   string       `json:"planned"`
	ETA       int          `json:"eta"` // in seconds
}

type TripArrival struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	State   ArrivalState `json:"state"`
	Planned string       `json:"planned"`
}

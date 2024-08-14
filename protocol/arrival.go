package protocol

type ArrivalState string

// State of an arrival
const (
	Undefined ArrivalState = ""
	Stopping  ArrivalState = "stopping"
	Predicted ArrivalState = "predicted"
	Planned   ArrivalState = "planned"
	Departed  ArrivalState = "departed"
)

// Type of a stop
type ArrivalType string

const (
	StopTypeBusStop     ArrivalType = "bus-stop"
	StopTypeParkingSpot ArrivalType = "parking-spot"
	StopTypeFerryStop   ArrivalType = "ferry-stop"
	StopTypeTrainStop   ArrivalType = "train-stop"
	StopTypeSubwayStop  ArrivalType = "subway-stop"
)

type StopArrival struct {
	Name      string       `json:"name"`
	Type      ArrivalType  `json:"type"`
	VehicleID string       `json:"vehicleId"`
	TripID    string       `json:"tripId"`
	RouteID   string       `json:"routeId"`
	RouteName string       `json:"routeName"`
	Direction string       `json:"direction"`
	State     ArrivalState `json:"state"`
	Planned   string       `json:"planned"`
	ETA       int          `json:"eta"` // in seconds
	Platform  string       `json:"platform"`
}

type TripArrival struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	State   ArrivalState `json:"state"`
	Planned string       `json:"planned"`
}

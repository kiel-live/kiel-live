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
type StopArrivalType string

const (
	StopTypeBusStop     StopArrivalType = "bus-stop"
	StopTypeParkingSpot StopArrivalType = "parking-spot"
	StopTypeFerryStop   StopArrivalType = "ferry-stop"
	StopTypeTrainStop   StopArrivalType = "train-stop"
	StopTypeSubwayStop  StopArrivalType = "subway-stop"
)

type StopArrival struct {
	Name      string          `json:"name"`
	Type      StopArrivalType `json:"type"`
	VehicleID string          `json:"vehicleId"`
	TripID    string          `json:"tripId"`
	RouteID   string          `json:"routeId"`
	RouteName string          `json:"routeName"`
	Direction string          `json:"direction"`
	State     ArrivalState    `json:"state"`
	Planned   string          `json:"planned"`
	ETA       int             `json:"eta"` // in seconds
	Platform  string          `json:"platform"`
}

type TripArrival struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	State   ArrivalState `json:"state"`
	Planned string       `json:"planned"`
}

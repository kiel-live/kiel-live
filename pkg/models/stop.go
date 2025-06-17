package models

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
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`   // Deprecated: use arrivals[].type or vehicles[].type instead
	Routes   []*Route       `json:"routes"` // list of routes using this stop
	Alerts   []string       `json:"alerts"` // general alerts for this stop
	Arrivals []*StopArrival `json:"arrivals"`
	Vehicles []*Vehicle     `json:"vehicles"`
	Location *Location      `json:"location"`
}

type ArrivalState string

// State of an arrival
const (
	Undefined ArrivalState = ""
	Stopping  ArrivalState = "stopping"
	Predicted ArrivalState = "predicted"
	Planned   ArrivalState = "planned"
	Departed  ArrivalState = "departed"
)

type StopArrival struct {
	Name      string       `json:"name"`
	Type      VehicleType  `json:"type"`
	VehicleID string       `json:"vehicleId"`
	TripID    string       `json:"tripId"`
	RouteID   string       `json:"routeId"`
	RouteName string       `json:"routeName"`
	Direction string       `json:"direction"`
	State     ArrivalState `json:"state"`
	Planned   string       `json:"planned"`
	Eta       int          `json:"eta"`
	Platform  string       `json:"platform"`
}

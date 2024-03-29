package protocol

// Type of a stop
type StopType string

const (
	StopTypeBusStop     StopType = "bus-stop"
	StopTypeParkingSpot StopType = "parking-spot"
	StopTypeFerryStop   StopType = "ferry-stop"
	StopTypeTrainStop   StopType = "train-stop"
	StopTypeSubwayStop  StopType = "subway-stop"
)

// Stop is a fixed point their for example a bus stop or a car-sharing parking spot is located.
type Stop struct {
	ID       string        `json:"id"`
	Provider string        `json:"provider"`
	Name     string        `json:"name"`
	Type     StopType      `json:"type"`
	Routes   []string      `json:"routes"` // list of routes using this stop
	Alerts   []string      `json:"alerts"` // general alerts for this stop
	Arrivals []StopArrival `json:"arrivals"`
	Location Location      `json:"location"`
}

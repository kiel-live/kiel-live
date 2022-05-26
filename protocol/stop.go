package protocol

// Type of a stop
const (
	StopTypeBusStop     = "bus-stop"
	StopTypeParkingSpot = "parking-spot"
	StopTypeFerryStop   = "ferry-stop"
	StopTypeTrainStop   = "train-stop"
	StopTypeSubwayStop  = "subway-stop"
)

// Stop is a fixed point their for example a bus stop or a car-sharing parking spot is located.
type Stop struct {
	ID       string        `json:"id"`
	Provider string        `json:"provider"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`   // use StopType...
	Routes   []string      `json:"routes"` // list of routes using this stop
	Alerts   []string      `json:"alerts"` // general alerts for this stop
	Arrivals []StopArrival `json:"arrivals"`
	Location Location      `json:"location"`
}

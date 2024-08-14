package protocol

// Stop is a fixed point their for example a bus stop or a car-sharing parking spot is located.
type Stop struct {
	ID       string        `json:"id"`
	Provider string        `json:"provider"`
	Name     string        `json:"name"`
	Type     ArrivalType   `json:"type"`   // TODO: deprecated, use arrivals[].type or vehicles[].type instead
	Routes   []string      `json:"routes"` // list of routes using this stop
	Alerts   []string      `json:"alerts"` // general alerts for this stop
	Arrivals []StopArrival `json:"arrivals"`
	Location Location      `json:"location"`
	Vehicles []Vehicle     `json:"vehicles"`
}

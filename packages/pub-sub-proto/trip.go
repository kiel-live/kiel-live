package protocol

// Trip is a tour represented by a list of stops executed by a vehicle (exp. bus) on a specific route.
type Trip struct {
	ID        string        `json:"id"`
	Provider  string        `json:"provider"`
	VehicleID string        `json:"vehicleId"`
	Direction string        `json:"direction"`
	Arrivals  []TripArrival `json:"arrivals"`
}

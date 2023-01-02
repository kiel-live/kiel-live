package protocol

// Type of a vehicle

type VehicleType string

const (
	VehicleTypeBus      VehicleType = "bus"
	VehicleTypeBike     VehicleType = "bike"
	VehicleTypeCar      VehicleType = "car"
	VehicleTypeEScooter VehicleType = "e-scooter"
	VehicleTypeFerry    VehicleType = "ferry"
	VehicleTypeTrain    VehicleType = "train"
	VehicleTypeSubway   VehicleType = "subway"
)

// Vehicle can be of a specific type (exp. bus, bike).
type Vehicle struct {
	ID       string      `json:"id"`
	Provider string      `json:"provider"`
	Name     string      `json:"name"`
	Type     VehicleType `json:"type"`
	State    string      `json:"state"`
	Battery  string      `json:"battery"` // in percent
	Location Location    `json:"location"`
	TripID   string      `json:"tripId"`
}

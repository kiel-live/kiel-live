package protocol

// Type of a vehicle
const (
	VehicleTypeBus      = "bus"
	VehicleTypeBike     = "bike"
	VehicleTypeCar      = "car"
	VehicleTypeEScooter = "e-scooter"
	VehicleTypeFerry    = "ferry"
	VehicleTypeTrain    = "train"
	VehicleTypeSubway   = "subway"
)

// Vehicle can be of a specific type (exp. bus, bike).
type Vehicle struct {
	ID       string   `json:"id"`
	Provider string   `json:"provider"`
	Name     string   `json:"name"`
	Type     string   `json:"type"` // use VehicleType...
	State    string   `json:"state"`
	Battery  string   `json:"battery"` // in percent
	Location Location `json:"location"`
	TripID   string   `json:"tripId"`
}

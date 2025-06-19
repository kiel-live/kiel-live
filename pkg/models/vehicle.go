package models

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
	VehicleTypeTram     VehicleType = "tram"
	VehicleTypeMoped    VehicleType = "moped"
	VehicleTypeEMoped   VehicleType = "e-moped"
)

// Vehicle is a moving object that can be of a specific type like bus or bike.
type Vehicle struct {
	ID          string      `json:"id"`
	Provider    string      `json:"provider"`
	Name        string      `json:"name"`
	Type        VehicleType `json:"type"`
	State       string      `json:"state"`
	Battery     string      `json:"battery"` // in percent
	Location    *Location   `json:"location"`
	TripID      string      `json:"tripId"`
	Actions     []*Action   `json:"actions"`
	Description string      `json:"description"`
}

func (v *Vehicle) String() string {
	if v == nil {
		return "Vehicle(nil)"
	}
	return "Vehicle(" + v.ID + ", " + v.Provider + ", " + v.Name + ")"
}

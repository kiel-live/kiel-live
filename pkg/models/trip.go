package models

// Trip is a tour represented by a list of stops executed by a vehicle (exp. bus) on a specific route.
type Trip struct {
	ID        string         `json:"id"`
	Provider  string         `json:"provider"`
	VehicleID string         `json:"vehicleId"`
	Direction string         `json:"direction"`
	Arrivals  []*TripArrival `json:"arrivals"`
	Path      []*Location    `json:"path"`
}

func (t *Trip) String() string {
	if t == nil {
		return "Trip(nil)"
	}
	return "Trip(" + t.ID + ", " + t.Provider + ", " + t.VehicleID + ")"
}

type TripArrival struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	State   ArrivalState `json:"state"`
	Planned string       `json:"planned"`
}

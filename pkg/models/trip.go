package models

// Trip is a tour represented by a list of stops executed by a vehicle (exp. bus) on a specific route.
type Trip struct {
	ID         string           `json:"id"`
	Provider   string           `json:"provider"`
	VehicleID  string           `json:"vehicleId"`
	Direction  string           `json:"direction"`
	Arrivals   []*TripArrival   `json:"arrivals,omitempty"`
	Departures []*TripDeparture `json:"departures,omitempty"`
	Path       []*Location      `json:"path,omitempty"`
}

func (t *Trip) String() string {
	if t == nil {
		return "Trip(nil)"
	}
	return "Trip(" + t.ID + ", " + t.Provider + ", " + t.VehicleID + ")"
}

type TripDeparture struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	State   DepartureState `json:"state"`
	Actual  string         `json:"actual"`
	Planned string         `json:"planned"`
}

// Deprecated: use TripDeparture instead
type TripArrival struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Planned string `json:"planned"`
}

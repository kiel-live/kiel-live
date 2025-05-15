package models

import "encoding/json"

type Trip struct {
	ID        string         `json:"id"`
	Provider  string         `json:"provider"`
	VehicleID string         `json:"vehicleId"`
	Direction string         `json:"direction"`
	Arrivals  []*TripArrival `json:"arrivals"`
	Path      []*Location    `json:"path"`
}

func (t *Trip) ToJSON() []byte {
	bytes, err := json.Marshal(t)
	if err != nil {
		return nil
	}
	return bytes
}

type TripArrival struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Planned string `json:"planned"`
}

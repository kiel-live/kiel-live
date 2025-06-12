package models

type Vehicle struct {
	ID       string    `json:"id"`
	Provider string    `json:"provider"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	State    string    `json:"state"`
	Battery  string    `json:"battery"`
	Location *Location `json:"location"`
	TripID   string    `json:"tripId"`
}

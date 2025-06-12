package models

type Stop struct {
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Routes   []*Route       `json:"routes"`
	Alerts   []string       `json:"alerts"`
	Arrivals []*StopArrival `json:"arrivals"`
	Vehicles []*Vehicle     `json:"vehicles"`
	Location *Location      `json:"location"`
}

type StopArrival struct {
	Name      string `json:"name"`
	VehicleID string `json:"vehicleId"`
	TripID    string `json:"tripId"`
	RouteID   string `json:"routeId"`
	RouteName string `json:"routeName"`
	Direction string `json:"direction"`
	State     string `json:"state"`
	Planned   string `json:"planned"`
	Eta       int    `json:"eta"`
	Platform  string `json:"platform"`
}

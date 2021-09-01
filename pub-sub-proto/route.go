package protocol

// Route is a fixed list of waypoints a vehicle could drive along. A one-time tour on a specific route is called a trip.
type Route struct {
	ID       string      `json:"id"`
	Provider string      `json:"provider"`
	Name     string      `json:"name"`
	Type     string      `json:"type"` // use VehicleType...
	IsActive bool        `json:"isActive"`
	Stops    []RouteStop `json:"stops"`
}

type RouteStop struct {
	ID       string   `json:"id"`
	Location Location `json:"location"`
}

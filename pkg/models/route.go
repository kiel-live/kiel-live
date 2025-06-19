package models

// Route is a fixed list of waypoints a vehicle could drive along. A one-time tour on a specific route is called a trip
type Route struct {
	ID       string       `json:"id"`
	Provider string       `json:"provider"`
	Name     string       `json:"name"`
	Type     string       `json:"type"` // TODO: use VehicleType?
	IsActive bool         `json:"isActive"`
	Stops    []*RouteStop `json:"stops"`
}

func (r *Route) String() string {
	if r == nil {
		return "Route(nil)"
	}
	return "Route(" + r.ID + ", " + r.Provider + ", " + r.Name + ")"
}

type RouteStop struct {
	ID       string    `json:"id"`
	Location *Location `json:"location"`
}

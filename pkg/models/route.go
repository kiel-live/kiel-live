package models

type Route struct {
	ID       string       `json:"id"`
	Provider string       `json:"provider"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	IsActive bool         `json:"isActive"`
	Stops    []*RouteStop `json:"stops"`
}

type RouteStop struct {
	ID       string    `json:"id"`
	Location *Location `json:"location"`
}

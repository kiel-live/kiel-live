package models

import "encoding/json"

type Route struct {
	ID       string       `json:"id"`
	Provider string       `json:"provider"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	IsActive bool         `json:"isActive"`
	Stops    []*RouteStop `json:"stops"`
}

func (r *Route) ToJSON() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}

type RouteStop struct {
	ID       string    `json:"id"`
	Location *Location `json:"location"`
}

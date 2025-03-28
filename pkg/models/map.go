package models

type Map struct {
	Stops    []*Stop    `json:"stops"`
	Vehicles []*Vehicle `json:"vehicles"`
}

package model

import (
	"encoding/json"

	"github.com/golang/geo/s2"
)

func (l *Location) GetCellID() s2.CellID {
	p := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
	return s2.CellIDFromLatLng(p).Parent(10)
}

func (l *LocationInput) ToLocation() *Location {
	return &Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Heading:   l.Heading,
	}
}

func (v *VehicleInput) ToVehicle() *Vehicle {
	return &Vehicle{
		ID:       v.ID,
		Provider: v.Provider,
		Name:     v.Name,
		Type:     v.Type,
		TripID:   v.TripID,
		State:    v.State,
		Battery:  v.Battery,
		Location: v.Location.ToLocation(),
	}
}

func (s *StopInput) ToStop() *Stop {
	return &Stop{
		ID:       s.ID,
		Provider: s.Provider,
		Name:     s.Name,
		Type:     s.Type,
		Location: s.Location.ToLocation(),
		// Routes: ,
		// Arrivals: ,
		// Alerts: ,
		// Vehicles: ,
	}
}

func (v *Vehicle) ToJSON() []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bytes
}

func (s *Stop) ToJSON() []byte {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return bytes
}

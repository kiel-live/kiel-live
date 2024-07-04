package model

import (
	"github.com/kiel-live/kiel-live/shared/models"
)

func (l *LocationInput) ToLocation() *models.Location {
	return &models.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Heading:   l.Heading,
	}
}

func (v *VehicleInput) ToVehicle() *models.Vehicle {
	return &models.Vehicle{
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

func (s *StopInput) ToStop() *models.Stop {
	return &models.Stop{
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

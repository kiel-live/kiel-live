package store

import "github.com/kiel-live/kiel-live/pkg/models"

func deriveStopMarker(s *models.Stop) *models.Stop {
	return &models.Stop{
		ID:       s.ID,
		Provider: s.Provider,
		Name:     s.Name,
		Type:     s.Type,
		Location: s.Location,
	}
}

func deriveVehicleMarker(v *models.Vehicle) *models.Vehicle {
	return &models.Vehicle{
		ID:       v.ID,
		Provider: v.Provider,
		Name:     v.Name,
		Type:     v.Type,
		State:    v.State,
		Location: v.Location,
	}
}

func locationEqual(a, b *models.Location) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Longitude != b.Longitude || a.Latitude != b.Latitude {
		return false
	}
	if a.Heading == nil && b.Heading == nil {
		return true
	}
	if a.Heading == nil || b.Heading == nil {
		return false
	}
	return *a.Heading == *b.Heading
}

func stopMarkersEqual(a, b *models.Stop) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ID == b.ID &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.Type == b.Type &&
		locationEqual(a.Location, b.Location)
}

func vehicleMarkersEqual(a, b *models.Vehicle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ID == b.ID &&
		a.Provider == b.Provider &&
		a.Name == b.Name &&
		a.Type == b.Type &&
		a.Description == b.Description &&
		a.Battery == b.Battery &&
		a.State == b.State &&
		locationEqual(a.Location, b.Location)
}

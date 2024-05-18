package model

import "github.com/golang/geo/s2"

func (l *Location) GetCellID() s2.CellID {
	p := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
	return s2.CellIDFromLatLng(p).Parent(11)
}

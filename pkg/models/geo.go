package models

import (
	"github.com/golang/geo/s2"
)

const (
	MinLevel             = 12
	MaxLevel             = 12
	BoundingBoxCellLimit = 20 // max number of cells in a bounding box
)

type Location struct {
	Longitude int  `json:"longitude"`         // exp: 54.306 * 3600000 = longitude
	Latitude  int  `json:"latitude"`          // exp: 10.149 * 3600000 = latitude
	Heading   *int `json:"heading,omitempty"` // in degree
	// Latitude  float64 `json:"latitude"`  // exp: 54.306
	// Longitude float64 `json:"longitude"` // exp: 10.149
}

func (l *Location) DistanceToMeters(location *Location) float64 {
	p1 := s2.LatLngFromDegrees(toDegreesFloat(l.Latitude), toDegreesFloat(l.Longitude))
	p2 := s2.LatLngFromDegrees(toDegreesFloat(location.Latitude), toDegreesFloat(location.Longitude))
	return p1.Distance(p2).Radians() * 6371000.0 // Earth radius in meters
}

func toDegreesFloat(value int) float64 {
	return float64(value) / 3600000.0
}

func (l *Location) GetCellID() s2.CellID {
	p := s2.LatLngFromDegrees(toDegreesFloat(l.Latitude), toDegreesFloat(l.Longitude))
	return s2.CellIDFromLatLng(p).Parent(MinLevel)
}

type BoundingBox struct {
	North float64
	East  float64
	South float64
	West  float64
}

func (b *BoundingBox) GetCellIDs() []s2.CellID {
	p1 := s2.LatLngFromDegrees(b.North, b.East)
	p2 := s2.LatLngFromDegrees(b.South, b.West)
	r := s2.RectFromLatLng(p1).AddPoint(p2)
	rc := s2.RegionCoverer{MinLevel: MinLevel, MaxLevel: MaxLevel, MaxCells: BoundingBoxCellLimit}
	return rc.Covering(r)
}

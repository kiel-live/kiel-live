package models

import (
	"github.com/golang/geo/s2"
)

const (
	MinLevel = 1
	MaxLevel = 15
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   *int    `json:"heading"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (l *Location) GetCellID() s2.CellID {
	p := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
	return s2.CellIDFromLatLng(p).Parent(10)
}

func (l *Location) GetCellIDs() []s2.CellID {
	cells := make([]s2.CellID, 0)
	p := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
	c := s2.CellIDFromLatLng(p)
	for i := min(c.Level(), MaxLevel); i >= MinLevel; i-- {
		cells = append(cells, c.Parent(i))
	}
	return cells
}

type BoundingBox struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
}

func (b *BoundingBox) GetCellIDs() []s2.CellID {
	p1 := s2.LatLngFromDegrees(b.MinLat, b.MinLng)
	p2 := s2.LatLngFromDegrees(b.MaxLat, b.MaxLng)
	r := s2.RectFromLatLng(p1).AddPoint(p2)
	rc := s2.RegionCoverer{MinLevel: MinLevel, MaxLevel: MaxLevel, MaxCells: 10}
	return rc.Covering(r)
}

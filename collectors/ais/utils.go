package main

import "math"

func calculateDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6371000.0

	lat1R := lat1 * math.Pi / 180
	lon1R := lon1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180
	lon2R := lon2 * math.Pi / 180

	dLat := lat2R - lat1R
	dLon := lon2R - lon1R

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1R)*math.Cos(lat2R)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}

func calculateBearing(lat1, lon1, lat2, lon2 float64) int {
	lat1R := lat1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180
	dLonR := (lon2 - lon1) * math.Pi / 180
	x := math.Sin(dLonR) * math.Cos(lat2R)
	y := math.Cos(lat1R)*math.Sin(lat2R) - math.Sin(lat1R)*math.Cos(lat2R)*math.Cos(dLonR)
	bearing := math.Atan2(x, y) * 180 / math.Pi
	return int(math.Mod(bearing+360, 360))
}

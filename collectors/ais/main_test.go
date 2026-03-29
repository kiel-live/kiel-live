package main

import (
	"math"
	"testing"
)

func TestCalculateDistanceMeters(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64
	}{
		{
			name: "same point",
			lat1: 54.3, lon1: 10.1,
			lat2: 54.3, lon2: 10.1,
			expected: 0,
			delta:    0.001,
		},
		{
			name: "short distance in Kiel (~1km)",
			lat1: 54.3233, lon1: 10.1228,
			lat2: 54.3323, lon2: 10.1228,
			expected: 1001,
			delta:    5,
		},
		{
			name: "equator 1 degree longitude (~111km)",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 1,
			expected: 111195,
			delta:    1,
		},
		{
			name: "north pole to equator (~10000km)",
			lat1: 90, lon1: 0,
			lat2: 0, lon2: 0,
			expected: 10007543,
			delta:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateDistanceMeters(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.expected) > tt.delta {
				t.Errorf("calculateDistanceMeters(%v, %v, %v, %v) = %v, want %v ± %v",
					tt.lat1, tt.lon1, tt.lat2, tt.lon2, got, tt.expected, tt.delta)
			}
		})
	}
}

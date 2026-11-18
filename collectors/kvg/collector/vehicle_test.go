package collector

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
)

func TestIsSameLocation(t *testing.T) {
	heading := func(h int) *int { return &h }

	tests := []struct {
		name string
		a    *models.Location
		b    *models.Location
		want bool
	}{
		{
			name: "same heading value in different allocations",
			a:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(90)},
			b:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(90)},
			want: true,
		},
		{
			name: "different heading",
			a:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(90)},
			b:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(180)},
			want: false,
		},
		{
			name: "both without heading",
			a:    &models.Location{Latitude: 1, Longitude: 2},
			b:    &models.Location{Latitude: 1, Longitude: 2},
			want: true,
		},
		{
			name: "only one with heading",
			a:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(0)},
			b:    &models.Location{Latitude: 1, Longitude: 2},
			want: false,
		},
		{
			name: "different position",
			a:    &models.Location{Latitude: 1, Longitude: 2, Heading: heading(90)},
			b:    &models.Location{Latitude: 3, Longitude: 2, Heading: heading(90)},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameLocation(tt.a, tt.b); got != tt.want {
				t.Errorf("isSameLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

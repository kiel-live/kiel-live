package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/models"
)

// --- modeToStopType ---

func TestModeToStopType(t *testing.T) {
	tests := []struct {
		modes []string
		want  models.StopType
	}{
		{[]string{"FERRY"}, models.StopTypeFerryStop},
		{[]string{"SUBWAY"}, models.StopTypeSubwayStop},
		{[]string{"REGIONAL_RAIL"}, models.StopTypeTrainStop},
		{[]string{"LONG_DISTANCE"}, models.StopTypeTrainStop},
		{[]string{"SUBURBAN"}, models.StopTypeTrainStop},
		{[]string{"TRAM"}, models.StopTypeBusStop}, // no tram stop type, falls through to bus
		{[]string{"BUS"}, models.StopTypeBusStop},
		{[]string{"REGIONAL_RAIL", "FERRY"}, models.StopTypeTrainStop}, // first match wins — REGIONAL_RAIL comes before FERRY
		{[]string{}, models.StopTypeBusStop},
	}

	for _, tt := range tests {
		got := modeToStopType(tt.modes)
		if got != tt.want {
			t.Errorf("modeToStopType(%v) = %q, want %q", tt.modes, got, tt.want)
		}
	}
}

// --- modeToVehicleType ---

func TestModeToVehicleType(t *testing.T) {
	tests := []struct {
		mode string
		want models.VehicleType
	}{
		{"TRAM", models.VehicleTypeTram},
		{"SUBWAY", models.VehicleTypeSubway},
		{"FERRY", models.VehicleTypeFerry},
		{"REGIONAL_RAIL", models.VehicleTypeTrain},
		{"LONG_DISTANCE", models.VehicleTypeTrain},
		{"SUBURBAN", models.VehicleTypeTrain},
		{"BUS", models.VehicleTypeBus},
		{"COACH", models.VehicleTypeBus},
		{"", models.VehicleTypeBus},
	}

	for _, tt := range tests {
		got := modeToVehicleType(tt.mode)
		if got != tt.want {
			t.Errorf("modeToVehicleType(%q) = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

// --- skipStop ---

func TestSkipStop(t *testing.T) {
	tests := []struct {
		name  string
		modes []string
		want  bool // true = should be skipped
	}{
		{"no modes", []string{}, true},
		{"bus only", []string{"BUS"}, true},
		{"coach only", []string{"COACH"}, true},
		{"tram", []string{"TRAM"}, false},
		{"ferry", []string{"FERRY"}, false},
		{"regional rail", []string{"REGIONAL_RAIL"}, false},
		{"bus and tram", []string{"BUS", "TRAM"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skipStop(place{Modes: tt.modes})
			if got != tt.want {
				t.Errorf("skipStop(modes=%v) = %v, want %v", tt.modes, got, tt.want)
			}
		})
	}
}

// --- mergeNearbyStops ---

// makeStop creates a minimal StopWithProviderID for merge tests.
// lat/lon are in degrees; they are stored as int×3600000.
func makeStop(id, name string, latDeg, lonDeg float64, providerID string) *StopWithProviderID {
	return &StopWithProviderID{
		Stop: &models.Stop{
			ID:   id,
			Name: name,
			Location: &models.Location{
				Latitude:  int(latDeg * 3600000),
				Longitude: int(lonDeg * 3600000),
			},
		},
		ProviderIDs: []string{providerID},
	}
}

func TestMergeNearbyStops(t *testing.T) {
	const (
		lat = 54.3145
		lon = 10.1305
	)

	t.Run("single stop passes through unchanged", func(t *testing.T) {
		stops := map[string]*StopWithProviderID{
			"motis-a": makeStop("motis-a", "Hbf", lat, lon, "de:1:a"),
		}
		got := mergeNearbyStops(stops)
		if len(got) != 1 {
			t.Fatalf("expected 1 stop, got %d", len(got))
		}
		if _, ok := got["motis-a"]; !ok {
			t.Error("expected motis-a to be present")
		}
	})

	t.Run("same name within 200m are merged", func(t *testing.T) {
		// +0.001 deg lat ≈ 111 m apart — well within 200 m
		stops := map[string]*StopWithProviderID{
			"motis-a": makeStop("motis-a", "Hbf", lat, lon, "de:1:a"),
			"motis-b": makeStop("motis-b", "Hbf", lat+0.001, lon, "de:1:b"),
		}
		got := mergeNearbyStops(stops)
		if len(got) != 1 {
			t.Fatalf("expected 1 stop after merge, got %d", len(got))
		}
	})

	t.Run("same name beyond 200m are kept separate", func(t *testing.T) {
		// +0.003 deg lat ≈ 333 m apart — outside 200 m radius
		stops := map[string]*StopWithProviderID{
			"motis-a": makeStop("motis-a", "Hbf", lat, lon, "de:1:a"),
			"motis-b": makeStop("motis-b", "Hbf", lat+0.003, lon, "de:1:b"),
		}
		got := mergeNearbyStops(stops)
		if len(got) != 2 {
			t.Fatalf("expected 2 stops to remain separate, got %d", len(got))
		}
	})

	t.Run("different names within 200m are kept separate", func(t *testing.T) {
		stops := map[string]*StopWithProviderID{
			"motis-a": makeStop("motis-a", "Hbf", lat, lon, "de:1:a"),
			"motis-b": makeStop("motis-b", "Westfriedhof", lat+0.001, lon, "de:1:b"),
		}
		got := mergeNearbyStops(stops)
		if len(got) != 2 {
			t.Fatalf("expected 2 stops (different names), got %d", len(got))
		}
	})

	t.Run("merged stop accumulates all provider IDs", func(t *testing.T) {
		stops := map[string]*StopWithProviderID{
			"motis-a": makeStop("motis-a", "Hbf", lat, lon, "de:1:a"),
			"motis-b": makeStop("motis-b", "Hbf", lat+0.001, lon, "de:1:b"),
		}
		got := mergeNearbyStops(stops)
		var canonical *StopWithProviderID
		for _, s := range got {
			canonical = s
		}
		if len(canonical.ProviderIDs) != 2 {
			t.Errorf("expected 2 provider IDs after merge, got %d: %v", len(canonical.ProviderIDs), canonical.ProviderIDs)
		}
	})

	// This is the regression test for the non-determinism bug: before the fix,
	// map iteration order was random, so "motis-z" could win over "motis-a"
	// on some runs, causing published stop IDs to flip between refreshes.
	t.Run("canonical ID is deterministic (smallest ID wins)", func(t *testing.T) {
		for i := range 100 {
			stops := map[string]*StopWithProviderID{
				"motis-z": makeStop("motis-z", "Hbf", lat, lon, "de:1:z"),
				"motis-a": makeStop("motis-a", "Hbf", lat+0.001, lon, "de:1:a"),
			}
			got := mergeNearbyStops(stops)
			if len(got) != 1 {
				t.Fatalf("iteration %d: expected 1 merged stop, got %d", i, len(got))
			}
			if _, ok := got["motis-a"]; !ok {
				t.Fatalf("iteration %d: motis-a should always be the canonical ID (it is lexicographically smaller than motis-z)", i)
			}
		}
	})
}

// --- GetStops (HTTP) ---

func TestGetStops(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/map/stops" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{ //nolint:errcheck
			{
				"name":   "Hauptbahnhof",
				"stopId": "de:01001:1",
				"lat":    54.3145,
				"lon":    10.1305,
				"modes":  []string{"REGIONAL_RAIL"},
			},
			// bus-only stop must be filtered out
			{
				"name":   "Bus only",
				"stopId": "de:01001:99",
				"lat":    54.32,
				"lon":    10.13,
				"modes":  []string{"BUS"},
			},
			// stop with no modes must be filtered out
			{
				"name":   "Unknown",
				"stopId": "de:01001:98",
				"lat":    54.31,
				"lon":    10.14,
				"modes":  []string{},
			},
		})
	}))
	defer srv.Close()

	t.Setenv("MOTIS_URL", srv.URL)

	stops, err := GetStops()
	if err != nil {
		t.Fatalf("GetStops() error: %v", err)
	}

	if len(stops) != 1 {
		t.Fatalf("expected 1 stop (bus-only and no-mode filtered), got %d", len(stops))
	}

	id := FormatMotisID("de:01001:1")
	s, ok := stops[id]
	if !ok {
		t.Fatalf("expected stop with ID %q to be present", id)
	}
	if s.Name != "Hauptbahnhof" {
		t.Errorf("Name = %q, want %q", s.Name, "Hauptbahnhof")
	}
	if s.Type != models.StopTypeTrainStop {
		t.Errorf("Type = %q, want %q", s.Type, models.StopTypeTrainStop)
	}
	if len(s.ProviderIDs) != 1 || s.ProviderIDs[0] != "de:01001:1" {
		t.Errorf("ProviderIDs = %v, want [de:01001:1]", s.ProviderIDs)
	}
}

func TestGetStopsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	t.Setenv("MOTIS_URL", srv.URL)

	_, err := GetStops()
	if err == nil {
		t.Error("expected error on HTTP 500, got nil")
	}
}

// --- GetStopDepartures (HTTP) ---

func TestGetStopDepartures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/stoptimes" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"stopTimes": []map[string]any{
				{
					"place": map[string]any{
						"departure":          "2024-06-01T14:30:00Z",
						"scheduledDeparture": "2024-06-01T14:28:00Z",
						"track":              "3",
					},
					"mode":           "REGIONAL_RAIL",
					"realTime":       true,
					"headsign":       "Hamburg Hbf",
					"tripId":         "trip:abc:1",
					"routeShortName": "RE70",
					"displayName":    "RE70",
					"cancelled":      false,
					"tripCancelled":  false,
				},
				// duplicate tripId must be deduped
				{
					"place": map[string]any{
						"departure":          "2024-06-01T14:30:00Z",
						"scheduledDeparture": "2024-06-01T14:28:00Z",
					},
					"mode":           "REGIONAL_RAIL",
					"realTime":       false,
					"headsign":       "Hamburg Hbf",
					"tripId":         "trip:abc:1",
					"routeShortName": "RE70",
					"cancelled":      false,
					"tripCancelled":  false,
				},
				// cancelled trip must be excluded
				{
					"place": map[string]any{
						"departure":          "2024-06-01T15:00:00Z",
						"scheduledDeparture": "2024-06-01T15:00:00Z",
					},
					"mode":          "REGIONAL_RAIL",
					"realTime":      false,
					"headsign":      "Flensburg",
					"tripId":        "trip:abc:2",
					"tripCancelled": true,
				},
			},
		})
	}))
	defer srv.Close()

	t.Setenv("MOTIS_URL", srv.URL)

	deps, err := GetStopDepartures([]string{"de:01001:1"})
	if err != nil {
		t.Fatalf("GetStopDepartures() error: %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("expected 1 departure (1 deduped, 1 trip-cancelled), got %d", len(deps))
	}

	d := deps[0]
	if d.RouteName != "RE70" {
		t.Errorf("RouteName = %q, want %q", d.RouteName, "RE70")
	}
	if d.Platform != "3" {
		t.Errorf("Platform = %q, want %q", d.Platform, "3")
	}
	if d.State != models.Predicted {
		t.Errorf("State = %q, want %q", d.State, models.Predicted)
	}
	if d.Type != models.VehicleTypeTrain {
		t.Errorf("Type = %q, want %q", d.Type, models.VehicleTypeTrain)
	}
}

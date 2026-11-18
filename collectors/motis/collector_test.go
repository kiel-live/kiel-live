package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/kiel-live/kiel-live/collectors/motis/api"
	"github.com/kiel-live/kiel-live/pkg/models"
)

// mockClient records UpdateStop / DeleteStop calls and satisfies client.Client.
type mockClient struct {
	mu           sync.Mutex
	updatedStops []*models.Stop
	deletedStops []string
}

func (m *mockClient) UpdateStop(stop *models.Stop) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *stop
	m.updatedStops = append(m.updatedStops, &cp)
	return nil
}

func (m *mockClient) DeleteStop(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deletedStops = append(m.deletedStops, id)
	return nil
}

// updatesFor returns all UpdateStop calls that targeted the given stop ID.
func (m *mockClient) updatesFor(id string) []*models.Stop {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []*models.Stop
	for _, s := range m.updatedStops {
		if s.ID == id {
			out = append(out, s)
		}
	}
	return out
}

// Unused interface methods.
func (m *mockClient) Connect() error                                          { return nil }
func (m *mockClient) Disconnect() error                                       { return nil }
func (m *mockClient) IsConnected() bool                                       { return true }
func (m *mockClient) SetOnConnectionChanged(fn func(bool))                    {}
func (m *mockClient) GetSubscribedTopics() []string                           { return nil }
func (m *mockClient) SetOnTopicsChanged(fn func(string, bool))                {}
func (m *mockClient) UpdateVehicle(v *models.Vehicle) error                   { return nil }
func (m *mockClient) UpdateTrip(tr *models.Trip) error                        { return nil }
func (m *mockClient) DeleteVehicle(id string) error                           { return nil }
func (m *mockClient) DeleteTrip(id string) error                              { return nil }

// newMockServer starts an httptest server that serves a single train stop and
// one departure for it. Tests call t.Setenv("MOTIS_URL", srv.URL) to wire it up.
func newMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/map/stops":
			json.NewEncoder(w).Encode([]map[string]any{ //nolint:errcheck
				{
					"name":   "Hauptbahnhof",
					"stopId": "de:01001:1",
					"lat":    54.3145,
					"lon":    10.1305,
					"modes":  []string{"REGIONAL_RAIL"},
				},
			})
		case "/api/v5/stoptimes":
			json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
				"stopTimes": []map[string]any{
					{
						"place": map[string]any{
							"departure":          "2024-06-01T14:30:00Z",
							"scheduledDeparture": "2024-06-01T14:28:00Z",
							"track":              "2",
						},
						"mode":           "REGIONAL_RAIL",
						"realTime":       true,
						"headsign":       "Hamburg Hbf",
						"tripId":         "trip:1",
						"routeShortName": "RE70",
						"cancelled":      false,
						"tripCancelled":  false,
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

const testStopMotisID = "de:01001:1"

var testStopID = api.FormatMotisID(testStopMotisID) // "motis-de-01001-1"

// --- subscribe ---

func TestSubscribeAddsStopToSubscribedSet(t *testing.T) {
	srv := newMockServer(t)
	t.Setenv("MOTIS_URL", srv.URL)

	mc := &mockClient{}
	c := &Collector{client: mc, subscribedStops: make(map[string]bool)}

	c.subscribe(testStopID)

	if !c.subscribedStops[testStopID] {
		t.Errorf("expected %q to be in subscribedStops after subscribe()", testStopID)
	}
}

func TestSubscribeImmediatelyPublishesDepartures(t *testing.T) {
	srv := newMockServer(t)
	t.Setenv("MOTIS_URL", srv.URL)

	mc := &mockClient{}
	c := &Collector{client: mc, subscribedStops: make(map[string]bool)}

	c.subscribe(testStopID)

	// At least one UpdateStop call for this stop should carry departures.
	updates := mc.updatesFor(testStopID)
	if len(updates) == 0 {
		t.Fatal("expected at least one UpdateStop call after subscribe(), got none")
	}
	last := updates[len(updates)-1]
	if last.Departures == nil {
		t.Error("expected Departures to be set after subscribe(), got nil")
	}
}

func TestSubscribeFetchesStopsWhenCacheIsEmpty(t *testing.T) {
	srv := newMockServer(t)
	t.Setenv("MOTIS_URL", srv.URL)

	mc := &mockClient{}
	// Collector with no cached stops (nil).
	c := &Collector{client: mc, subscribedStops: make(map[string]bool)}

	c.subscribe(testStopID)

	if c.stops == nil {
		t.Error("expected stops to be populated after subscribe() with empty cache")
	}
}

// --- unsubscribe ---

func TestUnsubscribeRemovesStopFromSubscribedSet(t *testing.T) {
	mc := &mockClient{}
	c := &Collector{
		client:          mc,
		subscribedStops: map[string]bool{testStopID: true},
	}

	c.unsubscribe(testStopID)

	if c.subscribedStops[testStopID] {
		t.Errorf("expected %q to be removed from subscribedStops after unsubscribe()", testStopID)
	}
}

func TestUnsubscribeUnknownStopIsNoop(t *testing.T) {
	mc := &mockClient{}
	c := &Collector{client: mc, subscribedStops: make(map[string]bool)}

	// Must not panic on unknown ID.
	c.unsubscribe("motis-does-not-exist")
}

// --- tick ---

func TestTickPublishesStopsWhenCacheExpired(t *testing.T) {
	srv := newMockServer(t)
	t.Setenv("MOTIS_URL", srv.URL)

	mc := &mockClient{}
	c := &Collector{
		client:          mc,
		subscribedStops: make(map[string]bool),
		// lastStopsFetch == 0 forces a refresh
	}

	c.tick()

	if len(mc.updatesFor(testStopID)) == 0 {
		t.Errorf("expected UpdateStop to be called for %q after tick() with expired cache", testStopID)
	}
}

func TestTickSkipsStopRefreshWhenCacheIsFresh(t *testing.T) {
	mc := &mockClient{}
	c := &Collector{
		client: mc,
		// Pre-populate stops so refreshDepartures can run without HTTP.
		stops: map[string]*api.StopWithProviderID{
			testStopID: {
				Stop:        &models.Stop{ID: testStopID, Name: "Hauptbahnhof"},
				ProviderIDs: []string{testStopMotisID},
			},
		},
		subscribedStops: make(map[string]bool),
		// Fresh cache — well within MaxCacheAge
		lastStopsFetch: time.Now().Unix(),
	}

	c.tick()

	// No HTTP server → no stop-list updates expected. If the cache check is
	// wrong this test would panic or produce an HTTP error log.
	if len(mc.deletedStops) != 0 {
		t.Errorf("expected no stops deleted on fresh cache tick, got: %v", mc.deletedStops)
	}
}

func TestTickRefreshesDeparturesForSubscribedStops(t *testing.T) {
	srv := newMockServer(t)
	t.Setenv("MOTIS_URL", srv.URL)

	mc := &mockClient{}
	c := &Collector{
		client: mc,
		stops: map[string]*api.StopWithProviderID{
			testStopID: {
				Stop:        &models.Stop{ID: testStopID, Name: "Hauptbahnhof"},
				ProviderIDs: []string{testStopMotisID},
			},
		},
		subscribedStops: map[string]bool{testStopID: true},
		lastStopsFetch:  time.Now().Unix(), // fresh — skip stop list refresh
	}

	c.tick()

	updates := mc.updatesFor(testStopID)
	if len(updates) == 0 {
		t.Fatal("expected UpdateStop call for subscribed stop after tick(), got none")
	}
	last := updates[len(updates)-1]
	if len(last.Departures) == 0 {
		t.Errorf("expected departures to be populated for subscribed stop after tick()")
	}
}

func TestTickDoesNotFetchDeparturesForUnsubscribedStops(t *testing.T) {
	// No HTTP server — any HTTP call would fail and surface in logs.
	mc := &mockClient{}
	c := &Collector{
		client: mc,
		stops: map[string]*api.StopWithProviderID{
			testStopID: {
				Stop:        &models.Stop{ID: testStopID, Name: "Hauptbahnhof"},
				ProviderIDs: []string{testStopMotisID},
			},
		},
		subscribedStops: make(map[string]bool), // nothing subscribed
		lastStopsFetch:  time.Now().Unix(),
	}

	c.tick()

	// No departure updates should have been issued.
	for _, s := range mc.updatesFor(testStopID) {
		if s.Departures != nil {
			t.Errorf("expected no departure update for unsubscribed stop, but got one with %d departures", len(s.Departures))
		}
	}
}

// --- reset ---

func TestResetClearsStopCache(t *testing.T) {
	mc := &mockClient{}
	c := &Collector{
		client: mc,
		stops: map[string]*api.StopWithProviderID{
			testStopID: {Stop: &models.Stop{ID: testStopID}},
		},
		subscribedStops: map[string]bool{testStopID: true},
		lastStopsFetch:  time.Now().Unix(),
	}

	c.reset()

	if c.stops != nil {
		t.Error("expected stops to be nil after reset()")
	}
	if c.lastStopsFetch != 0 {
		t.Error("expected lastStopsFetch to be 0 after reset()")
	}
}

func TestResetKeepsSubscribedStops(t *testing.T) {
	mc := &mockClient{}
	c := &Collector{
		client:          mc,
		subscribedStops: map[string]bool{testStopID: true},
		lastStopsFetch:  time.Now().Unix(),
	}

	c.reset()

	if !c.subscribedStops[testStopID] {
		t.Errorf("expected %q to remain in subscribedStops after reset()", testStopID)
	}
}

package store_test

import (
	"testing"
	"time"

	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// kieler viewport: roughly the city centre of Kiel.
var kielerViewport = &models.BoundingBox{
	North: 54.32, South: 54.30,
	East: 10.16, West: 10.10,
}

// intLoc converts decimal-degree lat/lon to the integer representation used by
// models.Location (degrees × 3 600 000).
func intLoc(lat, lon float64) *models.Location {
	return &models.Location{
		Latitude:  int(lat * 3_600_000),
		Longitude: int(lon * 3_600_000),
	}
}

func stopAt(id, name string, loc *models.Location) *models.Stop {
	return &models.Stop{ID: id, Provider: "test", Name: name, Type: "bus-stop", Location: loc}
}

// searchIDs is a helper that runs a search and returns the ordered stop IDs.
func searchIDs(s *store.Store, query string, bounds *models.BoundingBox) []string {
	results := s.Search(query, bounds, 20)
	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids
}

func makeLoc(lat, lon int) *models.Location {
	return &models.Location{Latitude: lat, Longitude: lon}
}

// subscribeCell subscribes ch to the map topic for the cell containing loc
// and returns the topic string.
func subscribeCell(s *store.Store, loc *models.Location, kind string, ch chan store.Event) string {
	cell := loc.GetCellID()
	topic := "map." + kind + "." + cell.ToToken()
	s.PubSub.Subscribe(topic, ch)
	return topic
}

func TestUpsertStopMarkerChanged(t *testing.T) {
	s := store.New()

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: makeLoc(195516960, 36539160)}
	changed := s.UpsertStop(stop)
	assert.True(t, changed, "first insert should report marker changed")

	// Same stop, only departures updated — marker unchanged
	stop2 := *stop
	stop2.Departures = []*models.StopDepartures{{Name: "Bus 1", Direction: "North"}}
	changed2 := s.UpsertStop(&stop2)
	assert.False(t, changed2, "departure-only update must not report marker changed")

	// Location moved
	stop3 := *stop
	stop3.Location = makeLoc(195517000, 36539160)
	changed3 := s.UpsertStop(&stop3)
	assert.True(t, changed3, "location change must report marker changed")
}

func TestMapTopicReceivesUpdates(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	ch := make(chan store.Event, 10)
	subscribeCell(s, loc, "stop", ch)

	vehicle := &models.Vehicle{
		ID: "ais-1", Provider: "ais", Name: "Ship", Type: "ferry", State: "running",
		Location: loc,
	}
	s.UpsertVehicle(vehicle)
	// vehicle is a different kind, stop channel should be empty
	assert.Empty(t, ch)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)
	require.Len(t, ch, 1)
	e := <-ch
	assert.Equal(t, store.EventActionUpdate, e.Action)
	assert.Equal(t, "stop", e.Kind)
	assert.Equal(t, "kvg-1", e.ID)
	assert.True(t, e.MarkerChanged)
}

func TestNoOpSuppression(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	ch := make(chan store.Event, 10)
	subscribeCell(s, loc, "vehicle", ch)

	vehicle := &models.Vehicle{
		ID: "ais-1", Provider: "ais", Name: "Ship", Type: "ferry", State: "running",
		Location: loc,
	}
	s.UpsertVehicle(vehicle)
	require.Len(t, ch, 1)
	e1 := <-ch
	assert.True(t, e1.MarkerChanged)

	// Same position again — marker unchanged, no map event
	s.UpsertVehicle(vehicle)
	assert.Empty(t, ch, "identical upsert must not publish to map topic")
}

func TestDetailTopicAlwaysReceivesUpdates(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	ch := make(chan store.Event, 10)
	s.PubSub.Subscribe("stop.kvg-1", ch)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)
	require.Len(t, ch, 1)
	e1 := <-ch
	assert.Equal(t, "stop.kvg-1", e1.Topic)
	assert.Equal(t, store.EventActionUpdate, e1.Action)

	// Departure-only update (marker unchanged) — detail topic still fires
	stop2 := *stop
	stop2.Departures = []*models.StopDepartures{{Name: "Bus 1", Direction: "North"}}
	s.UpsertStop(&stop2)
	require.Len(t, ch, 1)
	e2 := <-ch
	assert.Equal(t, "stop.kvg-1", e2.Topic)
	assert.False(t, e2.MarkerChanged)
	require.NotNil(t, e2.Entity)
	full := e2.Entity.(*models.Stop)
	require.Len(t, full.Departures, 1)
}

func TestCellFilterByKind(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)
	cell := loc.GetCellID()

	stopCh := make(chan store.Event, 10)
	vehicleCh := make(chan store.Event, 10)
	s.PubSub.Subscribe("map.stop."+cell.ToToken(), stopCh)
	s.PubSub.Subscribe("map.vehicle."+cell.ToToken(), vehicleCh)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)

	require.Len(t, stopCh, 1)
	assert.Empty(t, vehicleCh, "stop update must not reach vehicle channel")
}

func TestWrongCellGetsNothing(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	wrongCell := s2.CellIDFromLatLng(s2.LatLngFromDegrees(0, 0)).Parent(10)
	ch := make(chan store.Event, 10)
	s.PubSub.Subscribe("map.stop."+wrongCell.ToToken(), ch)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)

	assert.Empty(t, ch, "stop in a different cell must not reach wrong-cell subscriber")
}

func TestDeletePublishesToMapTopic(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	ch := make(chan store.Event, 10)
	subscribeCell(s, loc, "stop", ch)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)
	<-ch // consume upsert

	existed := s.DeleteStop("kvg-1")
	assert.True(t, existed)
	require.Len(t, ch, 1)
	e := <-ch
	assert.Equal(t, store.EventActionDelete, e.Action)
	assert.Equal(t, "kvg-1", e.ID)
}

func TestCellChangePublishesDeleteThenAdd(t *testing.T) {
	s := store.New()
	locA := makeLoc(195516960, 36539160)
	locB := makeLoc(100000000, 36539160) // far away — different level-10 cell

	cellA := locA.GetCellID()
	cellB := locB.GetCellID()
	require.NotEqual(t, cellA, cellB, "test requires two distinct cells")

	chA := make(chan store.Event, 10)
	chB := make(chan store.Event, 10)
	s.PubSub.Subscribe("map.stop."+cellA.ToToken(), chA)
	s.PubSub.Subscribe("map.stop."+cellB.ToToken(), chB)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: locA}
	s.UpsertStop(stop)
	require.Len(t, chA, 1)
	assert.Equal(t, store.EventActionUpdate, (<-chA).Action)
	assert.Empty(t, chB)

	// Move stop to cell B
	moved := *stop
	moved.Location = locB
	s.UpsertStop(&moved)

	// chA must receive a delete
	require.Len(t, chA, 1)
	assert.Equal(t, store.EventActionDelete, (<-chA).Action)

	// chB must receive an add
	require.Len(t, chB, 1)
	assert.Equal(t, store.EventActionUpdate, (<-chB).Action)
}

func TestGetMarkersInCells(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)
	cell := loc.GetCellID()
	wrongCell := s2.CellIDFromLatLng(s2.LatLngFromDegrees(0, 0)).Parent(10)

	stop := &models.Stop{ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop", Location: loc}
	s.UpsertStop(stop)

	stops, _ := s.GetMarkersInCells([]s2.CellID{cell})
	assert.Len(t, stops, 1)

	stops2, _ := s.GetMarkersInCells([]s2.CellID{wrongCell})
	assert.Empty(t, stops2)
}

func TestMarkerDerivation(t *testing.T) {
	s := store.New()
	loc := makeLoc(195516960, 36539160)

	ch := make(chan store.Event, 10)
	subscribeCell(s, loc, "stop", ch)

	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "A", Type: "bus-stop",
		Location:   loc,
		Departures: []*models.StopDepartures{{Name: "Bus 1"}},
		Alerts:     []string{"alert"},
	}
	s.UpsertStop(stop)

	require.Len(t, ch, 1)
	e := <-ch
	m := e.Marker.(*models.Stop)
	assert.Equal(t, stop.ID, m.ID)
	assert.Equal(t, stop.Provider, m.Provider)
	assert.Equal(t, stop.Name, m.Name)
	assert.Equal(t, stop.Type, m.Type)
	assert.Equal(t, stop.Location, m.Location)
	assert.Nil(t, m.Departures, "marker must not include departures")
	assert.Nil(t, m.Alerts, "marker must not include alerts")
}

// --- Search ordering tests ---

func TestSearch_NoBounds_BestTextMatchFirst(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("partial", "Hbf Vorplatz", intLoc(54.31, 10.13)))
	s.UpsertStop(stopAt("exact", "Hbf", intLoc(54.31, 10.13)))

	ids := searchIDs(s, "hbf", nil)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "exact", ids[0], "exact match must be first when no bounds given")
}

func TestSearch_WithBounds_InViewportPreferredForEqualScore(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("outside", "Kiel Hbf", intLoc(54.40, 10.20)))
	s.UpsertStop(stopAt("inside", "Kiel Ost", intLoc(54.31, 10.13)))

	ids := searchIDs(s, "kiel", kielerViewport)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "inside", ids[0], "in-viewport stop must be preferred for equal fuzzy scores")
}

func TestSearch_WithBounds_StrongTextMatchOutsideBeatsWeakMatchInside(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("weak-inside", "Hbf Lange Reihe Mitte", intLoc(54.31, 10.13)))
	s.UpsertStop(stopAt("strong-outside", "Hbf", intLoc(54.40, 10.20)))

	ids := searchIDs(s, "hbf", kielerViewport)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "strong-outside", ids[0], "clearly better text match must not be overridden by proximity boost")
}

func TestSearch_WithBounds_CloserToCenterRankedFirst(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("near-center", "Kiel Hbf", intLoc(54.31, 10.13)))
	s.UpsertStop(stopAt("near-edge", "Kiel Hbf", intLoc(54.32, 10.16)))

	ids := searchIDs(s, "kiel hbf", kielerViewport)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "near-center", ids[0], "stop closer to viewport center must rank higher")
}

func TestSearch_WithBounds_LongerInViewportBeatsShortOutside(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("dreikonen", "Dreikonen", intLoc(54.40, 10.20)))
	s.UpsertStop(stopAt("dreiecksplatz", "Dreiecksplatz", intLoc(54.31, 10.13)))

	ids := searchIDs(s, "drei", kielerViewport)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "dreiecksplatz", ids[0], "in-viewport stop must beat shorter outside stop")
}

func TestSearch_EmptyStore(t *testing.T) {
	s := store.New()
	assert.Empty(t, searchIDs(s, "kiel", kielerViewport))
}

func TestSearch_NilBoundsPreservesOrder(t *testing.T) {
	s := store.New()
	s.UpsertStop(stopAt("a", "Hauptbahnhof", intLoc(54.31, 10.13)))
	s.UpsertStop(stopAt("b", "Hauptbahnhof Vorplatz", intLoc(54.31, 10.13)))

	ids := searchIDs(s, "hauptbahnhof", nil)
	require.GreaterOrEqual(t, len(ids), 2)
	assert.Equal(t, "a", ids[0], "shorter/better match must be first without bounds")
}

// waitEvent reads one event from ch within 500ms or fails.
func waitEvent(t *testing.T, ch chan store.Event) store.Event {
	t.Helper()
	select {
	case e := <-ch:
		return e
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for store event")
		return store.Event{}
	}
}

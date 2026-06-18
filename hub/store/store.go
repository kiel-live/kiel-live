package store

import (
	"context"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/pkg/metrics"
	"github.com/kiel-live/kiel-live/pkg/pubsub"
	"github.com/kiel-live/kiel-live/pkg/search"
	"github.com/kiel-live/kiel-live/pkg/models"
)

const (
	StopTTL    = 30 * time.Minute
	TripTTL    = 5 * time.Minute
	VehicleTTL = 2 * time.Minute

	sweepInterval = 30 * time.Second
)

const EventActionUpdate = "update"
const EventActionDelete = "delete"

type Event struct {
	// Topic is the pub/sub topic this event was published to.
	// Map events use "map.{kind}.{cellToken}"; detail events use "{kind}.{id}".
	Topic         string
	Action        string    // "update" or "delete"
	Kind          string    // "stop", "vehicle", or "trip"
	ID            string
	CellID        s2.CellID // zero for trips or when location is nil
	MarkerChanged bool      // whether the sparse marker fields changed
	Marker        any       // *models.Stop (sparse) or *models.Vehicle (sparse); nil for trips/deletes
	Entity        any       // full *models.Stop, *models.Vehicle, or *models.Trip
}

type Store struct {
	stops    SpatialIndex[models.Stop, models.Stop]
	vehicles SpatialIndex[models.Vehicle, models.Vehicle]
	trips    BaseIndex[models.Trip, models.Trip]
	search   *search.Index

	// PubSub is the event bus. Map events are published to "map.{kind}.{cellToken}";
	// detail events to "{kind}.{id}". Callers subscribe before reading snapshots to
	// avoid missing concurrent updates.
	PubSub *pubsub.PubSub[Event]
}

func New() *Store {
	return &Store{
		stops:    newSpatialIndex(deriveStopMarker, stopMarkersEqual, func(s *models.Stop) *models.Location { return s.Location }, StopTTL),
		vehicles: newSpatialIndex(deriveVehicleMarker, vehicleMarkersEqual, func(v *models.Vehicle) *models.Location { return v.Location }, VehicleTTL),
		trips:    newBaseIndex(func(t *models.Trip) *models.Trip { return t }, func(a, b *models.Trip) bool { return a == b }, TripTTL),
		search:   search.New(),
		PubSub:   pubsub.New[Event](),
	}
}

// Start launches the TTL sweeper. It runs until ctx is cancelled.
func (s *Store) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(sweepInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.sweepExpired()
			}
		}
	}()
}

func (s *Store) sweepExpired() {
	for _, id := range s.stops.ExpiredIDs() {
		slog.Debug("stop ttl expired", "id", id)
		s.DeleteStop(id)
	}
	for _, id := range s.vehicles.ExpiredIDs() {
		slog.Debug("vehicle ttl expired", "id", id)
		s.DeleteVehicle(id)
	}
	for _, id := range s.trips.ExpiredIDs() {
		slog.Debug("trip ttl expired", "id", id)
		s.DeleteTrip(id)
	}
}

func (s *Store) UpsertStop(stop *models.Stop) bool {
	newMarker, newCell, oldCell, markerChanged := s.stops.Upsert(stop.ID, stop)
	s.search.Unregister(stop.ID)
	s.search.Register(stop.Name, stop.ID)

	if oldCell != 0 && oldCell != newCell {
		// Entity moved cells: notify old-cell subscribers to remove it.
		oldTopic := "map.stop." + oldCell.ToToken()
		s.PubSub.Publish(oldTopic, Event{
			Topic:  oldTopic,
			Action: EventActionDelete,
			Kind:   "stop",
			ID:     stop.ID,
			CellID: oldCell,
		})
	}

	if markerChanged || oldCell != newCell {
		mapTopic := "map.stop." + newCell.ToToken()
		s.PubSub.Publish(mapTopic, Event{
			Topic:         mapTopic,
			Action:        EventActionUpdate,
			Kind:          "stop",
			ID:            stop.ID,
			CellID:        newCell,
			MarkerChanged: markerChanged,
			Marker:        newMarker,
			Entity:        stop,
		})
	}

	// Always notify detail subscribers (departures etc. may have changed).
	detailTopic := "stop." + stop.ID
	s.PubSub.Publish(detailTopic, Event{
		Topic:  detailTopic,
		Action: EventActionUpdate,
		Kind:   "stop",
		ID:     stop.ID,
		Entity: stop,
	})

	return markerChanged
}

func (s *Store) DeleteStop(id string) bool {
	cell, existed := s.stops.Delete(id)
	if !existed {
		return false
	}
	s.search.Unregister(id)
	mapTopic := "map.stop." + cell.ToToken()
	s.PubSub.Publish(mapTopic, Event{
		Topic:  mapTopic,
		Action: EventActionDelete,
		Kind:   "stop",
		ID:     id,
		CellID: cell,
	})
	return true
}

func (s *Store) UpsertVehicle(vehicle *models.Vehicle) bool {
	newMarker, newCell, oldCell, markerChanged := s.vehicles.Upsert(vehicle.ID, vehicle)

	if oldCell != 0 && oldCell != newCell {
		oldTopic := "map.vehicle." + oldCell.ToToken()
		s.PubSub.Publish(oldTopic, Event{
			Topic:  oldTopic,
			Action: EventActionDelete,
			Kind:   "vehicle",
			ID:     vehicle.ID,
			CellID: oldCell,
		})
	}

	if markerChanged || oldCell != newCell {
		mapTopic := "map.vehicle." + newCell.ToToken()
		s.PubSub.Publish(mapTopic, Event{
			Topic:         mapTopic,
			Action:        EventActionUpdate,
			Kind:          "vehicle",
			ID:            vehicle.ID,
			CellID:        newCell,
			MarkerChanged: markerChanged,
			Marker:        newMarker,
			Entity:        vehicle,
		})
	}

	detailTopic := "vehicle." + vehicle.ID
	s.PubSub.Publish(detailTopic, Event{
		Topic:  detailTopic,
		Action: EventActionUpdate,
		Kind:   "vehicle",
		ID:     vehicle.ID,
		Entity: vehicle,
	})

	return markerChanged
}

func (s *Store) DeleteVehicle(id string) bool {
	cell, existed := s.vehicles.Delete(id)
	if !existed {
		return false
	}
	mapTopic := "map.vehicle." + cell.ToToken()
	s.PubSub.Publish(mapTopic, Event{
		Topic:  mapTopic,
		Action: EventActionDelete,
		Kind:   "vehicle",
		ID:     id,
		CellID: cell,
	})
	return true
}

func (s *Store) UpsertTrip(trip *models.Trip) {
	s.trips.Upsert(trip.ID, trip)
	detailTopic := "trip." + trip.ID
	s.PubSub.Publish(detailTopic, Event{
		Topic:  detailTopic,
		Action: EventActionUpdate,
		Kind:   "trip",
		ID:     trip.ID,
		Entity: trip,
	})
}

func (s *Store) DeleteTrip(id string) bool {
	if !s.trips.Delete(id) {
		return false
	}
	detailTopic := "trip." + id
	s.PubSub.Publish(detailTopic, Event{
		Topic:  detailTopic,
		Action: EventActionDelete,
		Kind:   "trip",
		ID:     id,
	})
	return true
}

func (s *Store) GetStop(id string) *models.Stop       { return s.stops.Get(id) }
func (s *Store) GetVehicle(id string) *models.Vehicle { return s.vehicles.Get(id) }
func (s *Store) GetTrip(id string) *models.Trip       { return s.trips.Get(id) }

// GetMarkersInCells returns current markers in the given cells.
// Passing nil cells returns all markers.
func (s *Store) GetMarkersInCells(cells []s2.CellID) ([]*models.Stop, []*models.Vehicle) {
	if cells == nil {
		return s.stops.AllMarkers(), s.vehicles.AllMarkers()
	}
	return s.stops.MarkersInCells(cells), s.vehicles.MarkersInCells(cells)
}

// maxProximityBoost is the maximum score multiplier awarded to a stop at the
// exact center of the viewport. Capped at 0.2 so proximity cannot override a
// meaningfully better text match.
const maxProximityBoost = 0.2

// Search returns up to limit stops that fuzzy-match query, ordered by a combined
// score of text relevance and proximity to the viewport center. When bounds is
// non-nil, in-viewport stops receive up to a 1.2× boost (highest at center).
func (s *Store) Search(query string, bounds *models.BoundingBox, limit int) []*models.Stop {
	hits := s.search.Search(query, limit)

	type entry struct {
		stop  *models.Stop
		score float64
	}
	entries := make([]entry, 0, len(hits))
	for _, h := range hits {
		stop := s.stops.Get(h.ID)
		if stop == nil {
			continue
		}
		score := h.Score * proximityMultiplier(stop.Location, bounds)
		entries = append(entries, entry{stop, score})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	results := make([]*models.Stop, len(entries))
	for i, e := range entries {
		results[i] = e.stop
	}
	return results
}

// proximityMultiplier returns a value in [1.0, 1+maxProximityBoost].
func proximityMultiplier(loc *models.Location, bounds *models.BoundingBox) float64 {
	if bounds == nil || loc == nil {
		return 1.0
	}
	lat := float64(loc.Latitude) / 3_600_000
	lng := float64(loc.Longitude) / 3_600_000

	if lat < bounds.South || lat > bounds.North || lng < bounds.West || lng > bounds.East {
		return 1.0
	}

	centerLat := (bounds.North + bounds.South) / 2
	centerLng := (bounds.East + bounds.West) / 2
	dlat := lat - centerLat
	dlng := lng - centerLng
	distSq := dlat*dlat + dlng*dlng

	halfH := (bounds.North - bounds.South) / 2
	halfW := (bounds.East - bounds.West) / 2
	halfDiagSq := halfH*halfH + halfW*halfW
	if halfDiagSq == 0 {
		return 1.0 + maxProximityBoost
	}

	proximity := 1.0 - math.Min(distSq/halfDiagSq, 1.0)
	return 1.0 + maxProximityBoost*proximity
}

func (s *Store) RegisterMetrics(reg *metrics.Registry) {
	reg.Register("store.stops", func() any { return s.stops.Len() })
	reg.Register("store.vehicles", func() any { return s.vehicles.Len() })
	reg.Register("store.trips", func() any { return s.trips.Len() })
}

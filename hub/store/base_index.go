package store

import (
	"sync"
	"time"
)

// BaseIndex is a generic entity store with a sparse marker view.
// E is the full entity type; M is the sparse marker type (often the same type with fewer fields set).
// Use SpatialIndex for entities with locations that need cell-based lookups.
type BaseIndex[E, M any] struct {
	deriveMarker func(*E) *M
	markersEqual func(*M, *M) bool
	ttl          time.Duration

	mu      sync.RWMutex
	items   map[string]*E
	markers map[string]*M // TODO: only use items and separate models into separate indexes
	expiry  map[string]time.Time
}

func newBaseIndex[E, M any](deriveMarker func(*E) *M, markersEqual func(*M, *M) bool, ttl time.Duration) BaseIndex[E, M] {
	return BaseIndex[E, M]{
		deriveMarker: deriveMarker,
		markersEqual: markersEqual,
		ttl:          ttl,
		items:        make(map[string]*E),
		markers:      make(map[string]*M),
		expiry:       make(map[string]time.Time),
	}
}

// Upsert stores item and its derived marker. Returns new marker and whether the marker changed.
func (idx *BaseIndex[E, M]) Upsert(id string, item *E) (newMarker *M, markerChanged bool) {
	newMarker = idx.deriveMarker(item)

	idx.mu.Lock()
	oldMarker := idx.markers[id]
	idx.items[id] = item
	idx.markers[id] = newMarker
	idx.expiry[id] = time.Now().Add(idx.ttl)
	idx.mu.Unlock()

	markerChanged = !idx.markersEqual(oldMarker, newMarker)
	return
}

// Delete removes the entity by id. Returns whether it existed.
func (idx *BaseIndex[E, M]) Delete(id string) (existed bool) {
	idx.mu.Lock()
	_, existed = idx.items[id]
	if existed {
		delete(idx.items, id)
		delete(idx.markers, id)
		delete(idx.expiry, id)
	}
	idx.mu.Unlock()
	return
}

// ExpiredIDs returns the IDs of all entities whose TTL has elapsed.
func (idx *BaseIndex[E, M]) ExpiredIDs() []string {
	now := time.Now()
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	var expired []string
	for id, t := range idx.expiry {
		if now.After(t) {
			expired = append(expired, id)
		}
	}
	return expired
}

// Get returns the full entity for id, or nil.
func (idx *BaseIndex[E, M]) Get(id string) *E {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.items[id]
}

// Len returns the number of stored entities.
func (idx *BaseIndex[E, M]) Len() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.items)
}

// AllMarkers returns all current sparse markers.
func (idx *BaseIndex[E, M]) AllMarkers() []*M {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]*M, 0, len(idx.markers))
	for _, m := range idx.markers {
		result = append(result, m)
	}
	return result
}

// AllItems returns all current full entities.
func (idx *BaseIndex[E, M]) AllItems() []*E {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]*E, 0, len(idx.items))
	for _, e := range idx.items {
		result = append(result, e)
	}
	return result
}

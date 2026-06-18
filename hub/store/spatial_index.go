package store

import (
	"sync"
	"time"

	"github.com/golang/geo/s2"
	"github.com/kiel-live/kiel-live/pkg/models"
)

// SpatialIndex embeds BaseIndex and adds an s2 cell index for location-based lookups.
type SpatialIndex[E, M any] struct {
	BaseIndex[E, M]
	getLocation func(*E) *models.Location

	// cellMu protects cellIndex and entityCell independently of the item/marker lock.
	cellMu     sync.RWMutex
	cellIndex  map[s2.CellID]map[string]struct{}
	entityCell map[string]s2.CellID
}

func newSpatialIndex[E, M any](
	deriveMarker func(*E) *M,
	markersEqual func(*M, *M) bool,
	getLocation func(*E) *models.Location,
	ttl time.Duration,
) SpatialIndex[E, M] {
	return SpatialIndex[E, M]{
		BaseIndex:   newBaseIndex(deriveMarker, markersEqual, ttl),
		getLocation: getLocation,
		cellIndex:   make(map[s2.CellID]map[string]struct{}),
		entityCell:  make(map[string]s2.CellID),
	}
}

// Upsert stores item and its derived marker, and updates the cell index.
// Returns new marker, new cell, old cell (zero if cell unchanged or entity is new), and whether the marker changed.
func (idx *SpatialIndex[E, M]) Upsert(id string, item *E) (newMarker *M, newCell s2.CellID, oldCell s2.CellID, markerChanged bool) {
	newMarker, markerChanged = idx.BaseIndex.Upsert(id, item)
	if loc := idx.getLocation(item); loc != nil {
		newCell = loc.GetCellID()
		oldCell = idx.setCell(id, newCell)
	}
	return
}

// Delete removes the entity and its cell index entry. Returns the cell it was in and whether it existed.
func (idx *SpatialIndex[E, M]) Delete(id string) (cell s2.CellID, existed bool) {
	existed = idx.BaseIndex.Delete(id)
	if existed {
		cell = idx.clearCell(id)
	}
	return
}

// MarkersInCells returns markers for entities in the given cells.
// Cell IDs are collected under cellMu first, then markers are looked up under RLock,
// to avoid holding both locks simultaneously.
func (idx *SpatialIndex[E, M]) MarkersInCells(cells []s2.CellID) []*M {
	idx.cellMu.RLock()
	var ids []string
	seen := make(map[string]struct{})
	for _, cell := range cells {
		for id := range idx.cellIndex[cell] {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				ids = append(ids, id)
			}
		}
	}
	idx.cellMu.RUnlock()

	idx.mu.RLock()
	result := make([]*M, 0, len(ids))
	for _, id := range ids {
		if m := idx.markers[id]; m != nil {
			result = append(result, m)
		}
	}
	idx.mu.RUnlock()
	return result
}

// setCell updates the cell for id and returns the previous cell.
// Returns zero for a brand-new entity (no previous cell); returns newCell itself
// when the entity existed but stayed in the same cell (so callers can distinguish
// "cell unchanged" from "new entity" by comparing oldCell == newCell).
func (idx *SpatialIndex[E, M]) setCell(id string, newCell s2.CellID) (oldCell s2.CellID) {
	idx.cellMu.Lock()
	defer idx.cellMu.Unlock()
	if prev, ok := idx.entityCell[id]; ok {
		oldCell = prev
		if prev == newCell {
			return // same cell: index unchanged, return actual old cell
		}
		if cells := idx.cellIndex[prev]; cells != nil {
			delete(cells, id)
			if len(cells) == 0 {
				delete(idx.cellIndex, prev)
			}
		}
	}
	idx.entityCell[id] = newCell
	if idx.cellIndex[newCell] == nil {
		idx.cellIndex[newCell] = make(map[string]struct{})
	}
	idx.cellIndex[newCell][id] = struct{}{}
	return oldCell
}

func (idx *SpatialIndex[E, M]) clearCell(id string) s2.CellID {
	idx.cellMu.Lock()
	defer idx.cellMu.Unlock()
	cell := idx.entityCell[id]
	if cells := idx.cellIndex[cell]; cells != nil {
		delete(cells, id)
		if len(cells) == 0 {
			delete(idx.cellIndex, cell)
		}
	}
	delete(idx.entityCell, id)
	return cell
}

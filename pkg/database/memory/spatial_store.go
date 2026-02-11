package memory

import (
	"sync"

	"github.com/golang/geo/s2"
)

// SpatialStore wraps EntityStore and adds spatial indexing
type SpatialStore[T any] struct {
	store      *Store[T]
	mu         sync.RWMutex
	cellsIndex *CellIndex
	getCellID  func(T) s2.CellID
	getID      func(T) string
}

// NewSpatialStore creates a new spatial entity store
func NewSpatialStore[T any](getID func(T) string, getCellID func(T) s2.CellID) *SpatialStore[T] {
	return &SpatialStore[T]{
		store:      NewStore(getID),
		cellsIndex: NewCellIndex(),
		getCellID:  getCellID,
		getID:      getID,
	}
}

// Get retrieves an entity by ID
func (s *SpatialStore[T]) Get(id string) (T, bool) {
	return s.store.Get(id)
}

// Set adds or updates an entity with spatial indexing
func (s *SpatialStore[T]) Set(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.getID(item)
	newCellID := s.getCellID(item)

	var oldCellID s2.CellID
	if oldItem, ok := s.store.Get(id); ok {
		oldCellID = s.getCellID(oldItem)
	}

	s.cellsIndex.UpdateItem(id, []s2.CellID{newCellID}, []s2.CellID{oldCellID})
	s.store.Set(item)
}

// Delete removes an entity and its spatial index entry
func (s *SpatialStore[T]) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item, ok := s.store.Get(id); ok {
		cellID := s.getCellID(item)
		s.cellsIndex.RemoveItem(id, []s2.CellID{cellID})
		return s.store.Delete(id)
	}

	return false
}

// GetInBounds returns all entities within the given cell IDs
func (s *SpatialStore[T]) GetInBounds(cellIDs []s2.CellID) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[string]bool)
	var results []T

	for _, cellID := range cellIDs {
		itemIDs := s.cellsIndex.GetItemIDs(cellID)
		for _, itemID := range itemIDs {
			if seen[itemID] {
				continue
			}
			seen[itemID] = true

			if item, ok := s.store.Get(itemID); ok {
				results = append(results, item)
			}
		}
	}

	return results
}

// Len returns the number of entities in the store
func (s *SpatialStore[T]) Len() int {
	return s.store.Len()
}

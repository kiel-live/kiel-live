package memory

import (
	"sync"
)

// Store is a generic store for entities
type Store[T any] struct {
	mu    sync.RWMutex
	items map[string]T
	getID func(T) string
}

// NewStore creates a new entity store
func NewStore[T any](getID func(T) string) *Store[T] {
	return &Store[T]{
		items: make(map[string]T),
		getID: getID,
	}
}

// Get retrieves an entity by ID
func (s *Store[T]) Get(id string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	return item, ok
}

// Set adds or updates an entity
func (s *Store[T]) Set(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.getID(item)
	s.items[id] = item
}

// Delete removes an entity
func (s *Store[T]) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[id]
	if !ok {
		return false
	}

	delete(s.items, id)
	return true
}

// Len returns the number of entities in the store
func (s *Store[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

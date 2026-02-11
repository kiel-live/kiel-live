package syncx

import "sync"

// MutexValue wraps a value of any type T and protects it with a mutex.
type MutexValue[T any] struct {
	mu  sync.Mutex
	val T
}

// NewMutexValue creates a new MutexValue with the given initial value.
func NewMutexValue[T any](initial T) *MutexValue[T] {
	return &MutexValue[T]{val: initial}
}

// Get safely retrieves a copy of the value.
func (m *MutexValue[T]) Get() T {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.val
}

// Set safely updates the value.
func (m *MutexValue[T]) Set(v T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.val = v
}

// Update safely updates the value using a provided function.
func (m *MutexValue[T]) Update(fn func(T) T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.val = fn(m.val)
}

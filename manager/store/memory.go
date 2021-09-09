package store

import (
	"errors"
	"sync"
)

type MemoryStore struct {
	store map[string]string
	sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Load() error {
	m.Lock()
	defer m.Unlock()
	m.store = make(map[string]string)
	return nil
}

func (m *MemoryStore) Get(key string) (string, error) {
	m.Lock()
	defer m.Unlock()
	value, ok := m.store[key]
	if !ok {
		return "", errors.New("value not found")
	}

	return value, nil
}

func (m *MemoryStore) Set(key string, value string) error {
	m.Lock()
	defer m.Unlock()
	m.store[key] = value
	return nil
}

func (m *MemoryStore) Delete(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.store, key)
	return nil
}

func (m *MemoryStore) Unload() error {
	m.Lock()
	defer m.Unlock()
	return nil
}

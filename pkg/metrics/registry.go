package metrics

import (
	"strings"
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	gauges map[string]func() any
}

func New() *Registry {
	return &Registry{gauges: make(map[string]func() any)}
}

func (r *Registry) Register(name string, fn func() any) {
	r.mu.Lock()
	r.gauges[name] = func() any { return fn() }
	r.mu.Unlock()
}

// Snapshot calls all registered gauge functions and returns a nested map
// built by splitting dot-notation keys (e.g. "store.stops" → {"store":{"stops":N}}).
func (r *Registry) Snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]any, len(r.gauges))
	for key, fn := range r.gauges {
		setNested(result, key, fn())
	}
	return result
}

func setNested(m map[string]any, key string, val any) {
	prefix, rest, found := strings.Cut(key, ".")
	if !found {
		m[key] = val
		return
	}
	sub, ok := m[prefix].(map[string]any)
	if !ok {
		sub = make(map[string]any)
		m[prefix] = sub
	}
	setNested(sub, rest, val)
}

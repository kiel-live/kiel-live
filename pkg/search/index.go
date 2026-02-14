package search

import (
	"strings"
	"sync"
)

// IndexConfig holds the configuration for a search index
type IndexConfig[T any] struct {
}

// SearchIndex is a generic full-text search index
type SearchIndex[T any] struct {
	mu sync.RWMutex

	getID      func(T) string
	getEntries func(T) []SearchableEntry

	// Index data structures
	terms     map[string]map[string]float64 // term -> id -> score
	itemTerms map[string][]string           // id -> terms
}

// NewSearchIndex creates a new generic search index
func NewSearchIndex[T any](
	getID func(T) string,
	getEntries func(T) []SearchableEntry,
) *SearchIndex[T] {
	return &SearchIndex[T]{
		getID:      getID,
		getEntries: getEntries,
		terms:      make(map[string]map[string]float64),
		itemTerms:  make(map[string][]string),
	}
}

// Set adds or updates an entity in the search index
func (si *SearchIndex[T]) Set(item T) {
	si.mu.Lock()
	defer si.mu.Unlock()

	id := si.getID(item)

	// Remove old terms if entity was already indexed
	si.removeTermsLocked(id)

	// Index new terms
	termsSet := make(map[string]float64)
	entries := si.getEntries(item)
	for _, entry := range entries {
		terms := Tokenize(entry.Text)
		for _, term := range terms {
			// Keep highest score for duplicate terms
			if existing, ok := termsSet[term]; !ok || entry.Score > existing {
				termsSet[term] = entry.Score
			}
		}
	}

	// Store terms and update index
	var allTerms []string
	for term, score := range termsSet {
		allTerms = append(allTerms, term)
		if si.terms[term] == nil {
			si.terms[term] = make(map[string]float64)
		}
		si.terms[term][id] = score
	}
	si.itemTerms[id] = allTerms
}

// Delete removes an entity from the search index by ID
func (si *SearchIndex[T]) Delete(id string) {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.removeTermsLocked(id)
}

// removeTermsLocked removes all terms associated with an ID
// Must be called with lock held
func (si *SearchIndex[T]) removeTermsLocked(id string) {
	if terms, exists := si.itemTerms[id]; exists {
		for _, term := range terms {
			if ids, ok := si.terms[term]; ok {
				delete(ids, id)
				if len(ids) == 0 {
					delete(si.terms, term)
				}
			}
		}
		delete(si.itemTerms, id)
	}
}

// Search performs a full-text search
func (si *SearchIndex[T]) Search(query string) []SearchResult {
	si.mu.RLock()
	defer si.mu.RUnlock()

	terms := Tokenize(query)
	if len(terms) == 0 {
		return []SearchResult{}
	}

	scores := make(map[string]float64)

	for _, term := range terms {
		// Exact matches
		if entities, ok := si.terms[term]; ok {
			for id, score := range entities {
				scores[id] += score
			}
		}

		// Prefix matches for autocomplete (only for terms >= 3 chars)
		if len(term) >= 3 {
			for indexTerm, entities := range si.terms {
				if strings.HasPrefix(indexTerm, term) && indexTerm != term {
					for id, score := range entities {
						scores[id] += score * 0.7
					}
				}
			}
		}
	}

	var results []SearchResult
	for id, score := range scores {
		results = append(results, SearchResult{
			ID:    id,
			Score: score,
		})
	}

	return results
}

// Len returns the number of entities in the index
func (si *SearchIndex[T]) Len() int {
	si.mu.RLock()
	defer si.mu.RUnlock()

	return len(si.itemTerms)
}

// Clear removes all entities from the index
func (si *SearchIndex[T]) Clear() {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.terms = make(map[string]map[string]float64)
	si.itemTerms = make(map[string][]string)
}

package search

import (
	"fmt"
	"sort"
	"sync"
)

// WeightedSearcher wraps a Searcher with a weight multiplier
type WeightedSearcher struct {
	Name     string
	Searcher Searcher
	Weight   float64
}

// MultiIndex combines multiple search indexes and returns unified results with configurable weights
type MultiIndex struct {
	mu      sync.RWMutex
	indexes []WeightedSearcher
}

// NewMultiIndex creates a new multi-index search
func NewMultiIndex() *MultiIndex {
	return &MultiIndex{
		indexes: make([]WeightedSearcher, 0),
	}
}

// AddIndex adds a search index to the multi-index with default weight of 1.0
func (mi *MultiIndex) AddIndex(name string, index Searcher) *MultiIndex {
	mi.AddWeightedIndex(name, index, 1.0)
	return mi
}

// AddWeightedIndex adds a search index with a specific weight multiplier
// Higher weights increase the relevance of results from that index
func (mi *MultiIndex) AddWeightedIndex(name string, index Searcher, weight float64) *MultiIndex {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	mi.indexes = append(mi.indexes, WeightedSearcher{
		Name:     name,
		Searcher: index,
		Weight:   weight,
	})

	return mi
}

// Search performs a weighted search across all indexes
func (mi *MultiIndex) Search(query string, limit int) []SearchResult {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	resultsMap := make(map[string]SearchResult)

	// search all indexes in parallel
	type indexResult struct {
		results []SearchResult
		name    string
		weight  float64
	}

	resultChans := make([]chan indexResult, len(mi.indexes))

	for i, ws := range mi.indexes {
		resultChans[i] = make(chan indexResult, 1)
		go func() {
			resultChans[i] <- indexResult{
				results: ws.Searcher.Search(query),
				name:    ws.Name,
				weight:  ws.Weight,
			}
		}()
	}

	// collect and merge results
	for _, ch := range resultChans {
		ir := <-ch
		for _, result := range ir.results {
			// create a unique key combining
			key := fmt.Sprintf("%s:%s", ir.name, result.ID)

			if existing, ok := resultsMap[key]; ok {
				// accumulate weighted scores for duplicates
				existing.Score += result.Score * ir.weight
				resultsMap[key] = existing
			} else {
				// apply weight to score
				multiResult := SearchResult{
					Index: ir.name,
					ID:    result.ID,
					Score: result.Score,
				}
				multiResult.Score *= ir.weight
				resultsMap[key] = multiResult
			}
		}
	}

	// convert map to slice
	var allResults []SearchResult
	for _, result := range resultsMap {
		allResults = append(allResults, result)
	}

	// sort by score descending
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Score > allResults[j].Score
	})

	// limit results
	if len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults
}

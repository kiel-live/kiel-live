package search_test

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/search"
)

// Example demonstrates the new search index API
func Test_basicUsage(t *testing.T) {
	// Create an index for strings
	index := search.NewSearchIndex(
		func(s string) string { return s }, // getID
		func(s string) []search.SearchableEntry { // getEntries
			return []search.SearchableEntry{
				{Text: s, Score: 1.0},
			}
		},
	)

	// Set some items and clear the index
	index.Set("Wien")
	index.Set("Graz")
	index.Clear()

	// Set items
	index.Set("Berlin")
	index.Set("Hamburg")
	index.Set("Munich")
	index.Set("Cologne")
	index.Set("Frankfurt")
	index.Set("Stuttgart")
	index.Set("Schweinfurt")

	// Update an item (just call Set again)
	index.Set("Berlin") // Updates existing entry

	// Delete by ID only - no need for the full object
	index.Delete("Hamburg")

	// Search
	results := index.Search("ber")
	if len(results) != 1 || results[0].ID != "Berlin" {
		t.Errorf("Expected to find Berlin, got %v", results)
	}

	// Search with no matches
	results = index.Search("xyz")
	if len(results) != 0 {
		t.Errorf("Expected no results, got %v", results)
	}

	// Search for deleted item
	results = index.Search("ham")
	if len(results) != 0 {
		t.Errorf("Expected no results for deleted item, got %v", results)
	}

	// Search for multiple matches
	results = index.Search("furt")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'furt', got %d", len(results))
	}

	// Search on updated index
	index.Set("Hamburg")
	results = index.Search("ham")
	if len(results) != 1 || results[0].ID != "Hamburg" {
		t.Errorf("Expected to find Hamburg after re-adding, got %v", results)
	}
}

// Example demonstrates weighted multi-index search
func Test_weightedSearch(t *testing.T) {
	// Create two indexes
	stopsIndex := search.NewSearchIndex(
		func(s string) string { return s },
		func(s string) []search.SearchableEntry {
			return []search.SearchableEntry{{Text: s, Score: 1.0}}
		},
	)

	vehiclesIndex := search.NewSearchIndex(
		func(s string) string { return s },
		func(s string) []search.SearchableEntry {
			return []search.SearchableEntry{{Text: s, Score: 1.0}}
		},
	)

	stopsIndex.Set("Central Station")
	stopsIndex.Set("Park Street")
	stopsIndex.Set("City Hall")
	vehiclesIndex.Set("Central Bus")

	// Create multi-index with weights
	// Stops get 2x weight, vehicles get 1x weight
	multi := search.NewMultiIndex()
	multi.AddWeightedIndex("stops", stopsIndex, 2.0)       // Prioritize stops
	multi.AddWeightedIndex("vehicles", vehiclesIndex, 1.0) // Lower priority for vehicles

	// Search across both indexes
	results := multi.Search("central", 10)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check that the stop is ranked higher than the vehicle due to weighting
	if results[0].ID != "Central Station" || results[0].Score <= results[1].Score {
		t.Errorf("Expected Central Station to be ranked higher than Central Bus, got %v", results)
	}
}

func Test_structuredData(t *testing.T) {
	type animal struct {
		Name string
		Type string
	}

	// Create an index for animals
	index := search.NewSearchIndex(
		func(a animal) string { return a.Name }, // getID
		func(a animal) []search.SearchableEntry { // getEntries
			return []search.SearchableEntry{
				{Text: a.Name, Score: 1.0},
				{Text: a.Type, Score: 0.8},
			}
		},
	)

	index.Set(animal{Name: "Fido", Type: "Dog"})
	index.Set(animal{Name: "Whiskers", Type: "Cat"})

	results := index.Search("dog")

	if len(results) != 1 || results[0].ID != "Fido" {
		t.Errorf("Expected to find Fido, got %v", results)
	}
}

package search

// SearchableEntry represents a piece of text with a relevance score
type SearchableEntry struct {
	Text  string
	Score float64
}

// SearchResult represents a search result item
type SearchResult struct {
	ID    string  `json:"id"`
	Index string  `json:"index"` // Only set for multi-index results
	Score float64 `json:"score"`
}

// Searcher is an interface for searchable indexes
type Searcher interface {
	Search(query string) []SearchResult
}

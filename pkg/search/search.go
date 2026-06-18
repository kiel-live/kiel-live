package search

import (
	"sort"
	"strings"
	"sync"
)

// Index is a fuzzy inverted index mapping search terms to entity IDs.
// Multiple terms can be registered per ID (e.g. stop name + route names).
type Index struct {
	mu    sync.RWMutex
	terms map[string]map[string]struct{} // normalized term → set of IDs
	byID  map[string]map[string]struct{} // ID → set of normalized terms (for cleanup)
}

func New() *Index {
	return &Index{
		terms: make(map[string]map[string]struct{}),
		byID:  make(map[string]map[string]struct{}),
	}
}

// Register indexes term for id. Multiple terms per ID are allowed.
func (idx *Index) Register(term, id string) {
	t := normalize(term)
	if t == "" {
		return
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if idx.terms[t] == nil {
		idx.terms[t] = make(map[string]struct{})
	}
	idx.terms[t][id] = struct{}{}
	if idx.byID[id] == nil {
		idx.byID[id] = make(map[string]struct{})
	}
	idx.byID[id][t] = struct{}{}
}

// Unregister removes all terms associated with id.
func (idx *Index) Unregister(id string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	for t := range idx.byID[id] {
		delete(idx.terms[t], id)
		if len(idx.terms[t]) == 0 {
			delete(idx.terms, t)
		}
	}
	delete(idx.byID, id)
}

type Result struct {
	ID    string
	Score float64
}

// Search returns up to limit results whose registered terms fuzzy-match query,
// ordered by descending match score.
func (idx *Index) Search(query string, limit int) []Result {
	q := normalize(query)
	if q == "" {
		return nil
	}

	// Score every registered term; keep best score per ID.
	scores := make(map[string]float64)

	idx.mu.RLock()
	for term, ids := range idx.terms {
		s := fuzzyScore(q, term)
		if s <= 0 {
			continue
		}
		for id := range ids {
			if s > scores[id] {
				scores[id] = s
			}
		}
	}
	idx.mu.RUnlock()

	ranked := make([]Result, 0, len(scores))
	for id, s := range scores {
		ranked = append(ranked, Result{ID: id, Score: s})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	return ranked
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// minSubstringScore is the floor for any substring match. Without this floor the
// length ratio alone would rank shorter-but-weaker names above longer-but-better
// names (e.g. "Dreikonen" over "Dreiecksplatz" for query "drei").
const minSubstringScore = 0.75

// fuzzyScore returns a score in (0, 1] for how well query matches term.
// Returns 0 if the match is too weak to be considered.
func fuzzyScore(query, term string) float64 {
	// Exact substring: strong match, score is the higher of the length-coverage
	// ratio and the minimum floor so proximity can act as a fair tiebreaker.
	if strings.Contains(term, query) {
		score := 0.5 + 0.5*(float64(len(query))/float64(len(term)))
		if score < minSubstringScore {
			score = minSubstringScore
		}
		return score
	}

	// Fall back to trigram Dice coefficient.
	qt := buildTrigrams(query)
	tt := buildTrigrams(term)
	if len(qt) == 0 || len(tt) == 0 {
		return 0
	}

	shared := 0
	for tg := range qt {
		if tt[tg] {
			shared++
		}
	}

	dice := float64(2*shared) / float64(len(qt)+len(tt))
	if dice < 0.3 {
		return 0
	}
	return dice
}

// buildTrigrams returns the set of character trigrams for s, padded with spaces.
func buildTrigrams(s string) map[string]bool {
	s = " " + s + " "
	tg := make(map[string]bool, len(s))
	for i := 0; i+3 <= len(s); i++ {
		tg[s[i:i+3]] = true
	}
	return tg
}

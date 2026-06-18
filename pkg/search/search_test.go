package search

import (
	"slices"
	"testing"
)

func ids(results []Result) []string {
	out := make([]string, len(results))
	for i, r := range results {
		out[i] = r.ID
	}
	return out
}

func TestIndex_ExactSubstring(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")
	idx.Register("Holstenplatz", "stop2")

	got := ids(idx.Search("haupt", 10))
	if !slices.Contains(got, "stop1") {
		t.Fatalf("expected stop1 in results, got %v", got)
	}
	if slices.Contains(got, "stop2") {
		t.Fatalf("stop2 should not match 'haupt', got %v", got)
	}
}

func TestIndex_CaseInsensitive(t *testing.T) {
	idx := New()
	idx.Register("Kiel Hauptbahnhof", "stop1")

	if got := ids(idx.Search("HAUPT", 10)); !slices.Contains(got, "stop1") {
		t.Fatalf("expected stop1 for uppercase query, got %v", got)
	}
}

func TestIndex_FuzzyTypo(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")

	got := ids(idx.Search("haupbahnhof", 10))
	if !slices.Contains(got, "stop1") {
		t.Fatalf("expected stop1 for typo query, got %v", got)
	}
}

func TestIndex_Unregister(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")
	idx.Unregister("stop1")

	if got := idx.Search("hauptbahnhof", 10); len(got) != 0 {
		t.Fatalf("expected no results after unregister, got %v", got)
	}
}

func TestIndex_MultipleTermsPerID(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")
	idx.Register("Kiel Hbf", "stop1")

	got := ids(idx.Search("hbf", 10))
	if !slices.Contains(got, "stop1") {
		t.Fatalf("expected stop1 via second term, got %v", got)
	}
}

func TestIndex_UnregisterClearsAllTerms(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")
	idx.Register("Kiel Hbf", "stop1")
	idx.Unregister("stop1")

	for _, q := range []string{"hauptbahnhof", "hbf"} {
		if got := idx.Search(q, 10); len(got) != 0 {
			t.Fatalf("query %q: expected empty after unregister, got %v", q, got)
		}
	}
}

func TestIndex_Limit(t *testing.T) {
	idx := New()
	for i := range 10 {
		idx.Register("Kiel Stop", string(rune('A'+i)))
	}

	got := idx.Search("kiel", 3)
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d: %v", len(got), got)
	}
}

func TestIndex_EmptyQuery(t *testing.T) {
	idx := New()
	idx.Register("Hauptbahnhof", "stop1")

	if got := idx.Search("", 10); len(got) != 0 {
		t.Fatalf("expected no results for empty query, got %v", got)
	}
}

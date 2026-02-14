package search

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Tokenize converts text into normalized search terms
func Tokenize(text string) []string {
	normalized := Normalize(text)
	normalized = strings.ToLower(normalized)

	var terms []string
	var current strings.Builder

	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				term := current.String()
				if len(term) >= 2 {
					terms = append(terms, term)
				}
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		term := current.String()
		if len(term) >= 2 {
			terms = append(terms, term)
		}
	}

	return terms
}

// Normalize removes accents and diacritics from text
func Normalize(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)
	return result
}

package shared

import (
	"strings"
	"unicode"
)

// isSeparator returns true for runes that should *split* tokens.
// - Letters → false  (stay inside the token)
// - Digits  → true   (acts as a boundary, so numbers are skipped)
// - Everything else (punctuation, whitespace, symbols) → true
func isSeparator(r rune) bool {
	return !unicode.IsLetter(r)
}

// Common stop words to filter out during indexing
var COMMON_WORDS map[string]bool = map[string]bool{
	"the":  true,
	"and":  true,
	"of":   true,
	"to":   true,
	"in":   true,
	"on":   true,
	"at":   true,
	"for":  true,
	"by":   true,
	"with": true,
	"as":   true,
	"from": true,
	"is":   true,
	"are":  true,
	"was":  true,
	"were": true,
	"be":   true,
	"been": true,
}

// Tokenize takes text and populates a posting map for the given document
// Converts to lowercase, removes punctuation, and filters out stop words
func Tokenize(text string) (map[string]int, map[string][]int) {
	termFreq := make(map[string]int)
	positions := make(map[string][]int)
	for i, word := range strings.FieldsFunc(text, isSeparator) {
		word = strings.Trim(strings.ToLower(word), ".,!?:;()[]\"'")
		if COMMON_WORDS[word] {
			continue
		}
		termFreq[word]++
		positions[word] = append(positions[word], i)
	}

	return termFreq, positions
}

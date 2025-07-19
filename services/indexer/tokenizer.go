package main

import (
	"strings"
)

// tokenise takes text and populates a posting map for the given document
// Converts to lowercase, removes punctuation, and filters out stop words
func tokenise(text string) map[string]int {
	termFreq := make(map[string]int)
	for _, word := range strings.FieldsFunc(text, isSeparator) {
		word = strings.Trim(strings.ToLower(word), ".,!?:;()[]\"'")
		if COMMON_WORDS[word] {
			continue
		}
		termFreq[word]++
	}

	return termFreq
}

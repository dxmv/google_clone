package main

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

// SearchResult represents a search result with score
type SearchResult struct {
	docMeta DocMeta
	score   float64
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, k int, db *badger.DB) []SearchResult {
	var results []SearchResult
	// parse and tokenize query
	queryTerms := tokenise(query) // will be a map of term to frequency

	fmt.Println("Query terms: ", queryTerms)

	// for each term, fetch postings
	var postingsForQuery []Posting
	for term, _ := range queryTerms {
		postingsForQuery = append(postingsForQuery, getPostings(db, term)...)
	}
	fmt.Println("Postings for query: ", postingsForQuery[0].DocID, postingsForQuery[0].Count)

	// score documents

	// rank and return top k results
	return results
}

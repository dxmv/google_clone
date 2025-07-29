package main

import (
	"fmt"
	"sort"

	badger "github.com/dgraph-io/badger/v4"
)

// SearchResult represents a search result with score
type SearchResult struct {
	DocMetadata DocMetadata `json:"docMeta"`
	Score       float64     `json:"score"`
	CountTerm   int         `json:"countTerm"`
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, db *badger.DB) []SearchResult {
	var results []SearchResult
	// parse and tokenize query
	queryTerms := tokenise(query) // will be a map of term to frequency

	fmt.Println("Query terms: ", queryTerms)

	// for each term, fetch postings
	docMap := make(map[string]SearchResult)
	for term, _ := range queryTerms {
		posting := getPostings(db, term)
		for _, posting := range posting {
			if _, ok := docMap[string(posting.DocID)]; !ok {
				docMap[string(posting.DocID)] = SearchResult{
					CountTerm:   0,
					Score:       0,
					DocMetadata: getMetadata(db, posting.DocID),
				}
			}
			docMap[string(posting.DocID)] = SearchResult{
				CountTerm:   docMap[string(posting.DocID)].CountTerm + 1,
				Score:       docMap[string(posting.DocID)].Score + float64(posting.Count),
				DocMetadata: docMap[string(posting.DocID)].DocMetadata,
			}
		}
	}

	// filter only docs that contain all terms
	for _, result := range docMap {
		if result.CountTerm == len(queryTerms) {
			results = append(results, result)
		}
	}

	// sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

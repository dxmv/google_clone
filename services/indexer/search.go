package main

import (
	"fmt"
	"sort"

	badger "github.com/dgraph-io/badger/v4"
)

// SearchResult represents a search result with score
type SearchResult struct {
	docMeta   DocMeta
	score     float64
	countTerm int
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, k int, db *badger.DB) []SearchResult {
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
					countTerm: 0,
					score:     0,
					docMeta:   getMeta(db, posting.DocID),
				}
			}
			docMap[string(posting.DocID)] = SearchResult{
				countTerm: docMap[string(posting.DocID)].countTerm + 1,
				score:     docMap[string(posting.DocID)].score + float64(posting.Count),
				docMeta:   docMap[string(posting.DocID)].docMeta,
			}
		}
	}

	// filter only docs that contain all terms
	for _, result := range docMap {
		if result.countTerm == len(queryTerms) {
			results = append(results, result)
		}
	}

	// sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// keep only top k results
	if len(results) > k {
		results = results[:k]
	}

	return results
}

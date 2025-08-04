package main

import (
	"fmt"
	"math"
	"sort"

	badger "github.com/dgraph-io/badger/v4"
)

// SearchResult represents a search result with score
type SearchResult struct {
	DocMetadata DocMetadata `json:"docMeta"`
	Score       float64     `json:"score"`
	CountTerm   int         `json:"countTerm"`
}

// BM25 parameters
var K = 1.2
var B = 0.75

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, db *badger.DB, avgDocLength float64, collectionSize int64) []SearchResult {
	// parse and tokenize query
	queryTerms := tokenise(query) // will be a map of term to frequency

	fmt.Println("Query terms: ", queryTerms)

	// for each term, fetch postings
	docMap := make(map[string]SearchResult)
	for term, _ := range queryTerms {
		posting := getPostings(db, term)
		numberOfDocumentsWithTerm := len(posting)
		idf := calculateIDF(numberOfDocumentsWithTerm, collectionSize)
		for _, posting := range posting {
			docID := string(posting.DocID)
			// get the metadata for the document
			metadata := getMetadata(db, posting.DocID)
			// Use the BM25 formula to calculate the score
			top, bottom := calculateTopBottom(posting, metadata, avgDocLength)
			score := idf * (top / bottom)
			_, ok := docMap[docID]
			// if the document is not in the map, add it
			if !ok {
				docMap[docID] = SearchResult{
					DocMetadata: metadata,
					Score:       score,
					CountTerm:   1,
				}
			}
			// if the document is in the map, update the score
			docMap[docID] = SearchResult{
				DocMetadata: metadata,
				Score:       docMap[docID].Score + score,
				CountTerm:   docMap[docID].CountTerm + 1,
			}
		}
	}

	results := make([]SearchResult, 0, len(docMap))
	for _, result := range docMap {
		results = append(results, result)
	}
	// sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func calculateTopBottom(posting Posting, metadata DocMetadata, avgDocLength float64) (float64, float64) {
	docLength := float64(metadata.ContentLength)
	// the percentage of the term in the document
	termFrequency := float64(posting.Count)
	top := termFrequency * (K + 1)
	bottom := termFrequency + K*(1-B+B*(docLength/avgDocLength))
	if bottom == 0 {
		bottom = 0.5
	}
	return top, bottom
}

func calculateIDF(numberOfDocumentsWithTerm int, collectionSize int64) float64 {
	top := float64(collectionSize) - float64(numberOfDocumentsWithTerm) + 0.5
	bottom := float64(numberOfDocumentsWithTerm) + 0.5
	if bottom == 0 {
		bottom = 0.5
	}
	return math.Log(top/bottom + 1)
}

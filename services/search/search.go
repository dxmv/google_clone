package main

import (
	"log"
	"math"
	"runtime"
	"sort"
	"sync"

	shared "github.com/dxmv/google_clone/shared"
)

// SearchResult represents a search result with score
type SearchResult struct {
	DocMetadata shared.DocMetadata `json:"docMeta"`
	Score       float64            `json:"score"`
	CountTerm   int                `json:"countTerm"`
}

// BM25 parameters
var K = 1.2
var B = 0.75

type Job struct {
	Posting shared.Posting
	IDF     float64
}

func worker(id int, jobs <-chan Job, results chan<- SearchResult, wg *sync.WaitGroup, storage *shared.Storage, avgDocLength float64) {
	defer wg.Done()
	for job := range jobs {
		posting := job.Posting
		idf := job.IDF
		docID := string(posting.DocID)
		// get the metadata for the document
		metadata, err := storage.GetMetadata(docID)
		if err != nil {
			log.Println("Error getting metadata: ", err)
			continue
		}
		// Use the BM25 formula to calculate the score
		top, bottom := calculateTopBottom(posting, metadata, avgDocLength)
		score := idf * (top / bottom)
		results <- SearchResult{
			DocMetadata: metadata,
			Score:       score,
			CountTerm:   1,
		}
	}
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, storage *shared.Storage, avgDocLength float64, collectionSize int64) []SearchResult {
	// parse and tokenize query
	queryTerms := shared.Tokenize(query) // will be a map of term to frequency

	// for each term, fetch postings
	docMap := make(map[string]SearchResult)
	for term, _ := range queryTerms {
		posting := storage.GetPostings(term)
		numberOfDocumentsWithTerm := len(posting)
		idf := calculateIDF(numberOfDocumentsWithTerm, collectionSize)
		jobs := make(chan Job, len(posting))
		results := make(chan SearchResult, len(posting))
		wg := sync.WaitGroup{}

		workers := runtime.NumCPU()
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go worker(i, jobs, results, &wg, storage, avgDocLength)
		}

		go func() {
			for _, posting := range posting {
				jobs <- Job{Posting: posting, IDF: idf}
			}
			close(jobs)
		}()

		wg.Wait()
		close(results)

		// process the results
		for result := range results {
			docID := result.DocMetadata.Hash
			_, ok := docMap[docID]
			// if the document is not in the map, add it
			if !ok {
				docMap[docID] = SearchResult{
					DocMetadata: result.DocMetadata,
					Score:       result.Score,
					CountTerm:   1,
				}
			} else {
				// if the document is in the map, update the score
				docMap[docID] = SearchResult{
					DocMetadata: result.DocMetadata,
					Score:       docMap[docID].Score + result.Score,
					CountTerm:   docMap[docID].CountTerm + 1,
				}
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

func calculateTopBottom(posting shared.Posting, metadata shared.DocMetadata, avgDocLength float64) (float64, float64) {
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

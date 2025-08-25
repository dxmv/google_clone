package main

import (
	"container/heap"
	"log"
	"math"
	"runtime"
	"sync"

	shared "github.com/dxmv/google_clone/shared"
)

// SearchResult represents a search result with score
type SearchResult struct {
	Hash      string  `json:"hash"`
	Score     float64 `json:"score"`
	CountTerm int     `json:"countTerm"`
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
		docLength, err := storage.GetDocLength(docID)
		if err != nil {
			log.Println("Error getting metadata: ", err)
			continue
		}
		// Use the BM25 formula to calculate the score
		top, bottom := calculateTopBottom(posting, docLength, avgDocLength)
		score := idf * (top / bottom)
		results <- SearchResult{
			Hash:      docID,
			Score:     score,
			CountTerm: 1,
		}
	}
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, storage *shared.Storage, avgDocLength float64, collectionSize int64, cache *LRUCache[string, []SearchResult]) []SearchResult {
	// parse and tokenize query
	queryTerms := shared.Tokenize(query) // will be a map of term to frequency

	res, ok := cache.Get(query)
	if ok {
		return res
	}

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
			docID := result.Hash
			_, ok := docMap[docID]
			// if the document is not in the map, add it
			if !ok {
				docMap[docID] = SearchResult{
					Hash:      result.Hash,
					Score:     result.Score,
					CountTerm: 1,
				}
			} else {
				// if the document is in the map, update the score
				docMap[docID] = SearchResult{
					Hash:      result.Hash,
					Score:     docMap[docID].Score + result.Score,
					CountTerm: docMap[docID].CountTerm + 1,
				}
			}
		}
	}

	// Use SearchHeap to sort results by score (max-heap)
	searchHeap := &SearchHeap{}
	heap.Init(searchHeap)

	for _, result := range docMap {
		heap.Push(searchHeap, result)
	}

	// Convert heap to sorted slice (highest scores first)
	results := make([]SearchResult, 0, len(docMap))
	for searchHeap.Len() > 0 {
		result := heap.Pop(searchHeap).(SearchResult)
		results = append(results, result)
	}

	// cache the result
	cache.Put(query, results)

	return results
}

// searchPaginated performs search query and returns only the requested page of results
// More efficient for large result sets as it only sorts/returns what's needed
func searchPaginated(query string, storage *shared.Storage, avgDocLength float64, collectionSize int64, cache *LRUCache[string, []SearchResult], page, count int32) []SearchResult {
	// First check if we have the full results cached
	fullResults, ok := cache.Get(query)
	if ok {
		// If cached, just slice and return the requested page
		offset := (page - 1) * count
		if offset >= int32(len(fullResults)) {
			return []SearchResult{}
		}
		if offset+count > int32(len(fullResults)) {
			return fullResults[offset:]
		}
		return fullResults[offset : offset+count]
	}

	// If not cached, perform full search and cache it
	allResults := search(query, storage, avgDocLength, collectionSize, cache)

	// Return the requested page
	offset := (page - 1) * count
	if offset >= int32(len(allResults)) {
		return []SearchResult{}
	}
	if offset+count > int32(len(allResults)) {
		return allResults[offset:]
	}
	return allResults[offset : offset+count]
}

func calculateTopBottom(posting shared.Posting, docLength uint32, avgDocLength float64) (float64, float64) {
	// the percentage of the term in the document
	termFrequency := float64(posting.Count)
	top := termFrequency * (K + 1)
	res := float64(docLength) / avgDocLength
	bottom := termFrequency + K*(1-B+B*res)
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

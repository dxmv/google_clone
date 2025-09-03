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
var K = 2.0
var B = 0.9

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
	queryTerms, _ := shared.Tokenize(query) // will be a map of term to frequency

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

	// convert map keys to slice for indexing
	terms := make([]string, 0, len(queryTerms))
	for term := range queryTerms {
		terms = append(terms, term)
	}

	// Apply position scoring to all results
	for docID, result := range docMap {
		// calculate position bonuses
		if result.CountTerm >= 2 {
			bonus := 0.0

			// for each adjacent term
			for i := 0; i+1 < len(queryTerms); i++ {
				term1 := terms[i]
				term2 := terms[i+1]
				// get the position of the term in the document
				pos1 := storage.GetPositions(term1, result.Hash)
				pos2 := storage.GetPositions(term2, result.Hash)
				if len(pos1) == 0 || len(pos2) == 0 {
					continue
				}
				// check if phrases are next to each other
				hits := phraseHits(pos1, pos2)
				if hits > 0 {
					// limit hits to 3
					if hits > 3 {
						hits = 3
					}
					bonus += 0.5 * float64(hits)
				}
			}
			result.Score += bonus
			docMap[docID] = result // update the result in the map
		}
	}

	// Convert map to unsorted slice
	results := make([]SearchResult, 0, len(docMap))
	for _, result := range docMap {
		results = append(results, result)
	}

	// cache the unsorted results
	cache.Put(query, results)

	return results
}

// searchPaginated performs search query and returns only the requested page of results
// More efficient for large result sets as it only sorts/returns what's needed
func searchPaginated(query string, storage *shared.Storage, avgDocLength float64, collectionSize int64, cache *LRUCache[string, []SearchResult], page, count int32) ([]SearchResult, int64) {
	// Get unsorted results from search
	allResults := search(query, storage, avgDocLength, collectionSize, cache)
	totalResults := int64(len(allResults))

	// Calculate how many results we need (page * count)
	neededResults := page * count
	if neededResults > int32(len(allResults)) {
		neededResults = int32(len(allResults))
	}

	// Use heap to get top neededResults
	searchHeap := &SearchHeap{}
	heap.Init(searchHeap)

	// Add all results to heap
	for _, result := range allResults {
		heap.Push(searchHeap, result)
	}

	// Extract top neededResults from heap
	sortedResults := make([]SearchResult, 0, neededResults)
	for i := int32(0); i < neededResults && searchHeap.Len() > 0; i++ {
		result := heap.Pop(searchHeap).(SearchResult)
		sortedResults = append(sortedResults, result)
	}

	// Return the requested page
	offset := (page - 1) * count
	if offset >= int32(len(sortedResults)) {
		return []SearchResult{}, totalResults
	}
	if offset+count > int32(len(sortedResults)) {
		return sortedResults[offset:], totalResults
	}
	return sortedResults[offset : offset+count], totalResults
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

// check how many times phrases are to each other
func phraseHits(a, b []int) int {
	i, j, hits := 0, 0, 0
	for i < len(a) && j < len(b) {
		d := b[j] - a[i]
		if d == 1 {
			hits++
			i++
			j++
		} else if d > 1 {
			i++
		} else {
			j++
		}
	}
	return hits
}

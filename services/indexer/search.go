package main

// SearchResult represents a search result with score
type SearchResult struct {
	docMeta DocMeta
	score   float64
}

// search performs search query and returns top k results
// Fetches postings, scores documents, and ranks by relevance
func search(query string, k int) []SearchResult {
	// TODO: Implementation for search functionality
	// 1. Parse and tokenize query
	// 2. Fetch postings for query terms
	// 3. Score documents using TF-IDF or similar
	// 4. Rank and return top k results
	var results []SearchResult
	return results
}

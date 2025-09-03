package main

// An SearchHeap is a min-heap of SearchResult.
type SearchHeap []SearchResult

func (h SearchHeap) Len() int { return len(h) }
func (h SearchHeap) Less(i, j int) bool {
	if h[i].Score == h[j].Score {
		return h[i].CountTerm > h[j].CountTerm
	}
	return h[i].Score > h[j].Score
}
func (h SearchHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *SearchHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(SearchResult))
}

func (h *SearchHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

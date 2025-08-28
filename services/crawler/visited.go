package main

import "sync"

type Visited struct {
	mu      sync.Mutex
	visited map[string]bool
}

func NewVisited() *Visited {
	return &Visited{
		visited: make(map[string]bool),
		mu:      sync.Mutex{},
	}
}

func (v *Visited) Add(url string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.visited[url] = true
}

func (v *Visited) IsVisited(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.visited[url]
}

func (v *Visited) GetVisited() map[string]bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.visited
}

func (v *Visited) CheckAndMark(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.visited[url] {
		return true
	}

	v.visited[url] = true
	return false
}

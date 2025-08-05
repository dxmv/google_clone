package main

import "sync"

type SkippedJobs struct {
	mu      sync.Mutex
	skipped []Job
}

func NewSkippedJobs() *SkippedJobs {
	return &SkippedJobs{
		skipped: make([]Job, 0),
		mu:      sync.Mutex{},
	}
}

func (s *SkippedJobs) Add(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.skipped = append(s.skipped, job)
}

func (s *SkippedJobs) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.skipped)
}

func (s *SkippedJobs) GetAllAndClear() []Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make a copy of the jobs
	result := make([]Job, len(s.skipped))
	copy(result, s.skipped)

	// Clear the slice
	s.skipped = s.skipped[:0]

	return result
}

func (s *SkippedJobs) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.skipped) == 0
}

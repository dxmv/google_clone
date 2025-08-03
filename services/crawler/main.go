package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
)

var START_LINKS = []string{
	"https://en.wikipedia.org/wiki/Philosophy",
}

const MAX_DEPTH = 2

type Job struct {
	URL   string
	Depth int
}

type DocMetadata struct {
	URL   string
	Depth int
	Title string
	Hash  string
	Links []string
}

type Visited struct {
	mu      sync.Mutex
	visited map[string]bool
}

type SkippedJobs struct {
	mu      sync.Mutex
	skipped []Job
}

func processJob(job Job, jobs chan Job, skippedJobs *SkippedJobs, visited *Visited, wg *sync.WaitGroup) {
	// check if already visited
	visited.mu.Lock()
	if visited.visited[job.URL] {
		visited.mu.Unlock()
		return
	}
	// mark as visited
	visited.visited[job.URL] = true
	visited.mu.Unlock()

	docMetadata := DocMetadata{
		URL:   job.URL,
		Depth: job.Depth,
		Title: "",
		Hash:  "",
	}
	// get the html
	body, err := fetch(job.URL)
	if err != nil {
		fmt.Println("Error getting HTML from", job.URL, err)
		return
	}

	// extract the links from the html
	links := extractLinks(body, &docMetadata)
	docMetadata.Links = links
	// add new jobs to the queue
	for _, link := range links {
		newJob := Job{URL: link, Depth: job.Depth + 1}
		// check if depth is too high
		if newJob.Depth > MAX_DEPTH {
			continue
		}
		select {
		case jobs <- newJob:
			wg.Add(1)
		default:
			skippedJobs.mu.Lock()
			skippedJobs.skipped = append(skippedJobs.skipped, newJob)
			skippedJobs.mu.Unlock()
		}
	}

	// save the html
	hash := sha256.Sum256(body)
	err = saveHTML(hex.EncodeToString(hash[:]), body)
	if err != nil {
		fmt.Println("Error saving HTML", err)
	}

	// save metadata
	docMetadata.Hash = hex.EncodeToString(hash[:])
	err = saveMetadata(docMetadata)
	if err != nil {
		fmt.Println("Error saving metadata", err)
	}

}

func worker(id int, jobs chan Job, skippedJobs *SkippedJobs, visited *Visited, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Println("Worker", id, "processing job", job.URL, "depth", job.Depth)
		processJob(job, jobs, skippedJobs, visited, wg)
		wg.Done()
	}
}

func main() {
	createDirectory(PAGES_DIR)
	createDirectory(METADATA_DIR)

	// create visited map
	visited := Visited{
		visited: make(map[string]bool),
		mu:      sync.Mutex{},
	}

	// create jobs channel
	jobs := make(chan Job, 1000)
	skippedJobs := SkippedJobs{
		skipped: make([]Job, 0),
		mu:      sync.Mutex{},
	}
	wg := sync.WaitGroup{}

	// start workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(i, jobs, &skippedJobs, &visited, &wg)
	}

	// seed jobs
	for _, link := range START_LINKS {
		wg.Add(1)
		jobs <- Job{URL: link, Depth: 0}
	}

	wg.Wait()
	// process skipped jobs in 10 rounds
	rounds := 0
	for rounds < 10 {
		fmt.Printf("Waiting for round %d jobs to complete...\n", rounds+1)
		wg.Wait()

		skippedJobs.mu.Lock()
		if len(skippedJobs.skipped) == 0 {
			skippedJobs.mu.Unlock()
			fmt.Println("No more skipped jobs, finishing...")
			break
		}

		currentSkipped := make([]Job, len(skippedJobs.skipped))
		copy(currentSkipped, skippedJobs.skipped)
		skippedJobs.skipped = skippedJobs.skipped[:0] // Clear the slice
		skippedJobs.mu.Unlock()

		rounds++
		fmt.Printf("Round %d: Processing %d skipped jobs\n", rounds, len(currentSkipped))

		processed := 0
		for _, job := range currentSkipped {

			// Try to add to jobs
			select {
			case jobs <- job:
				wg.Add(1)
				processed++
			default:
				// Still can't send, add back to skipped
				skippedJobs.mu.Lock()
				skippedJobs.skipped = append(skippedJobs.skipped, job)
				skippedJobs.mu.Unlock()
			}
		}
		fmt.Printf("Round %d: Successfully queued %d jobs\n", rounds, processed)
	}

	// wait for all jobs to be processed & close jobs channel
	wg.Wait()
	close(jobs)

	// print visited URLs
	fmt.Println("\n\n--------------------------------")
	fmt.Println("Visited URLs:")
	for url := range visited.visited {
		fmt.Println(url)
	}
	fmt.Println(len(visited.visited))
	fmt.Println(len(skippedJobs.skipped))
}

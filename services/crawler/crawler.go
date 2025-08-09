package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type Job struct {
	URL   string
	Depth int
}

type DocMetadata struct {
	URL           string
	Depth         int
	Title         string
	Hash          string
	Links         []string
	ContentLength int
}

/*
storage: interface to save the html and metadata
config: config to configure the crawler
wg: wait group to wait for all workers to finish
jobs: channel to send jobs to the workers
skippedJobs: struct to store jobs that were skipped
visited: struct to store visited urls
*/
type Crawler struct {
	storage     Storage
	config      *Config
	wg          *sync.WaitGroup
	jobs        chan Job
	skippedJobs SkippedJobs
	visited     *Visited
}

// NewCrawler creates a new crawler with the given storage and config
func NewCrawler(storage Storage, config *Config) *Crawler {
	return &Crawler{
		storage: storage,
		config:  config,
		wg:      &sync.WaitGroup{},
		jobs:    make(chan Job, config.JobsBuffer),
		skippedJobs: SkippedJobs{
			skipped: make([]Job, 0),
		},
		visited: NewVisited(),
	}
}

// Start starts the crawler
func (c *Crawler) Start() error {
	c.storage.CreateHTMLDirectory(c.config.PagesDir)
	c.storage.CreateMetadataDirectory(c.config.MetadataDir)

	// Start workers
	for i := 0; i < c.config.NumWorkers; i++ {
		go c.worker(i)
	}

	// Seed initial jobs
	for _, link := range c.config.StartLinks {
		c.wg.Add(1)
		c.jobs <- Job{URL: link, Depth: 0}
	}

	// wait for all jobs to be processed
	c.wg.Wait()

	// process skipped jobs
	c.processSkippedJobs()

	// wait for all remaining jobs to be processed
	c.wg.Wait()
	close(c.jobs)

	// print results
	c.printResults()

	return nil
}

// worker is a worker that processes jobs
func (c *Crawler) worker(id int) {
	for job := range c.jobs {
		fmt.Println("Worker", id, "processing job", job.URL, "depth", job.Depth)
		c.processJob(job)
		c.wg.Done()
	}
}

// processJob processes a job
func (c *Crawler) processJob(job Job) {
	// check if already visited
	if c.visited.IsVisited(job.URL) {
		return
	}
	// mark as visited
	c.visited.Add(job.URL)

	docMetadata := DocMetadata{
		URL:           job.URL,
		Depth:         job.Depth,
		Title:         "",
		Hash:          "",
		ContentLength: 0,
		Links:         []string{},
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
	docMetadata.ContentLength = len(body)
	// add new jobs to the queue
	for _, link := range links {
		newJob := Job{URL: link, Depth: job.Depth + 1}
		// check if depth is too high
		if newJob.Depth > c.config.MaxDepth {
			continue
		}

		select {
		case c.jobs <- newJob:
			c.wg.Add(1)
		default:
			c.skippedJobs.Add(newJob)
		}
	}
	hash := sha256.Sum256(body)
	hashString := hex.EncodeToString(hash[:])

	// save the html
	err = c.storage.SaveHTML(hashString, body)
	if err != nil {
		fmt.Println("Error saving HTML", err)
	}

	// save metadata
	docMetadata.Hash = hashString
	err = c.storage.SaveMetadata(docMetadata)
	if err != nil {
		fmt.Println("Error saving metadata", err)
	}
}

// processSkippedJobs processes skipped jobs
func (c *Crawler) processSkippedJobs() {
	rounds := 0
	for rounds < c.config.MaxRounds {
		fmt.Printf("Waiting for round %d jobs to complete...\n", rounds+1)
		c.wg.Wait()

		if c.skippedJobs.IsEmpty() {
			fmt.Println("No more skipped jobs, finishing...")
			break
		}

		prevSkippedJobs := c.skippedJobs.GetAllAndClear()
		rounds++
		fmt.Printf("Round %d: Processing %d skipped jobs\n", rounds, len(prevSkippedJobs))

		processed := 0
		for _, job := range prevSkippedJobs {
			// Try to add to jobs
			select {
			case c.jobs <- job:
				c.wg.Add(1)
				processed++
			default:
				// Still can't send, add back to skipped
				c.skippedJobs.Add(job)
			}
		}
		fmt.Printf("Round %d: Successfully queued %d jobs\n", rounds, processed)
	}
}

// printResults prints the results of the crawler
func (c *Crawler) printResults() {
	fmt.Println("\n\n--------------------------------")
	fmt.Println("Visited URLs:", len(c.visited.GetVisited()))
	fmt.Println("Skipped jobs:", c.skippedJobs.Count())
}

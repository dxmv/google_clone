package main

import (
	"context"
	"log"
	"runtime"
	"sync"

	shared "github.com/dxmv/google_clone/shared"
)

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

type WorkerResult struct {
	Metadata shared.DocMetadata
	Postings map[string][]shared.Posting
}

func worker(id int, jobs <-chan shared.DocMetadata, results chan<- WorkerResult, wg *sync.WaitGroup, corpus shared.Corpus) {
	defer wg.Done()
	for j := range jobs {
		log.Println("worker", id, "processing job", j.Title)
		postings := make(map[string][]shared.Posting)
		// read the metadata file
		metadata := j
		// index the html file
		html, err := corpus.GetHTML(context.Background(), metadata.Hash+".html")
		error_check(err)
		index_file(string(html), []byte(metadata.Hash), postings)
		results <- WorkerResult{Metadata: metadata, Postings: postings}
	}
}

func makeIndex(storage *shared.Storage) {
	log.Println("Indexing...")
	// read the metadata directory
	docs, err := storage.ListMetadata()
	error_check(err)

	// create the jobs and results channels
	jobs := make(chan shared.DocMetadata, 100)
	results := make(chan WorkerResult, 100)
	wg := sync.WaitGroup{}

	// start the workers with the number of CPUs
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg, storage.Corpus)
	}

	// enqueue the jobs from the metadata directory
	go func() {
		for _, doc := range docs {
			jobs <- doc
		}
		close(jobs)
	}()

	// wait for the workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// merge the postings into a single map
	// merge the metadata into a single slice
	postings := make(map[string][]shared.Posting)
	metadata := make([]shared.DocMetadata, 0)
	stats := shared.Stats{
		AvgDocLength: 0,
		TotalDocs:    0,
	}
	for r := range results {
		for term, posting := range r.Postings {
			postings[term] = append(postings[term], posting...)
		}
		metadata = append(metadata, r.Metadata)
		stats.TotalDocs++
		stats.AvgDocLength += float64(r.Metadata.ContentLength)
	}
	if stats.TotalDocs > 0 {
		stats.AvgDocLength /= float64(stats.TotalDocs)
	}
	// save the postings to the database, term by term
	log.Println("Saving postings...", len(postings))
	for term, posting := range postings {
		log.Println("Saving postings for", term)
		err = storage.SavePostings(map[string][]shared.Posting{term: posting})
		error_check(err)
	}
	log.Println("Saving postings complete")

	// save the stats to the database
	log.Println("Saving stats...")
	err = storage.SaveStats(stats)
	error_check(err)
	log.Println("Saving stats complete")

	log.Println("Indexing complete")
}

func main() {
	corpus := shared.Corpus(shared.NewMinoMongoCorpus())

	storage := shared.NewStorage(corpus)
	defer storage.DB.Close()
	makeIndex(storage)
}

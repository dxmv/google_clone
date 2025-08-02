package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	badger "github.com/dgraph-io/badger/v4"
)

type DocMetadata struct {
	URL   string
	Depth int
	Title string
	Hash  string
	Links []string
}

const PAGES_DIR = "../crawler/pages"
const METADATA_DIR = "../crawler/metadata"

var (
	reindex = flag.Bool("reindex", false, "Rebuild the index before serving")
)

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

type WorkerResult struct {
	Metadata DocMetadata
	Postings map[string][]Posting
}

func worker(id int, jobs <-chan string, results chan<- WorkerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		log.Println("worker", id, "processing job", j)
		postings := make(map[string][]Posting)
		// read the metadata file
		metadataFile, err := os.ReadFile(j)
		error_check(err)
		// unmarshal the metadata file
		var metadata DocMetadata
		err = json.Unmarshal(metadataFile, &metadata)
		error_check(err)
		// index the html file
		htmlFilePath := filepath.Join(PAGES_DIR, metadata.Hash+".html")
		index_file(htmlFilePath, []byte(metadata.Hash), postings)
		results <- WorkerResult{Metadata: metadata, Postings: postings}
	}
}

func makeIndex(db *badger.DB) {
	log.Println("Indexing...")
	// read the metadata directory
	files, err := os.ReadDir(METADATA_DIR)
	error_check(err)

	// create the jobs and results channels
	jobs := make(chan string, 100)
	results := make(chan WorkerResult, 100)
	wg := sync.WaitGroup{}

	// start the workers with the number of CPUs
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// enqueue the jobs from the metadata directory
	go func() {
		for _, file := range files {
			jobs <- filepath.Join(METADATA_DIR, file.Name())
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
	postings := make(map[string][]Posting)
	metadata := make([]DocMetadata, 0)
	for r := range results {
		for term, posting := range r.Postings {
			postings[term] = append(postings[term], posting...)
		}
		metadata = append(metadata, r.Metadata)
	}
	// save the postings to the database, term by term
	log.Println("Saving postings...")
	for term, posting := range postings {
		log.Println("Saving postings for", term)
		err = savePostings(db, map[string][]Posting{term: posting})
		error_check(err)
	}
	log.Println("Saving postings complete")

	// save the metadata to the database
	log.Println("Saving metadata...")
	for _, m := range metadata {
		// save the metadata to the database
		log.Println("Saving metadata for", m.Title, "with hash", m.Hash)
		err = saveMetadata(db, []byte(m.Hash), m)
		error_check(err)
	}
	log.Println("Saving metadata complete")
	log.Println("Indexing complete")
}

func main() {
	flag.Parse()

	db, err := openDB()
	error_check(err)
	defer db.Close()
	if *reindex {
		makeIndex(db)
	}
	startServer(db)
}

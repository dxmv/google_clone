package main

import (
	"encoding/json"
	"fmt"
	"os"

	"flag"

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

// Global inverted index structure
var postings map[string][]Posting

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func makeIndex(db *badger.DB) {
	// Initialize the postings map
	postings = make(map[string][]Posting)

	// Read the metadata directory
	files, err := os.ReadDir(METADATA_DIR)
	error_check(err)

	// Index each file
	for _, file := range files {
		// open the metadata file
		metadataFilePath := fmt.Sprintf("%s/%s", METADATA_DIR, file.Name())
		metadataFile, err := os.ReadFile(metadataFilePath)
		error_check(err)
		var metadata DocMetadata
		err = json.Unmarshal(metadataFile, &metadata)
		error_check(err)

		// index the html file
		hash := metadata.Hash
		htmlFilePath := fmt.Sprintf("%s/%s", PAGES_DIR, hash+".html")
		index_file(htmlFilePath, []byte(hash), postings)
		fmt.Println("Indexed: ", metadata.Title, "with hash: ", hash)
		// save the metadata to the Badger database
		err = saveMetadata(db, []byte(hash), metadata)
		error_check(err)

	}

	fmt.Printf("Total terms indexed: %d\n", len(postings))

	for term, postingsList := range postings {
		singleTermMap := map[string][]Posting{term: postingsList}
		err = savePostings(db, singleTermMap)
		error_check(err)
	}

	if err != nil {
		fmt.Println("Error saving postings: ", err)
	} else {
		fmt.Println("Saved postings...")
	}
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

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"
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

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Global inverted index structure
var postings map[string][]Posting

func main() {
	// Initialize the postings map
	postings = make(map[string][]Posting)

	// Read the metadata directory
	files, err := os.ReadDir(METADATA_DIR)
	error_check(err)

	// Open the Badger database
	db, err := openDB()
	error_check(err)
	defer db.Close()

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

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		results := search(query, db)
		json.NewEncoder(w).Encode(results)
		fmt.Println("Results: ", results)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

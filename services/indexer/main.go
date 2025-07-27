package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"log"
	"net/http"
)

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
	// TODO: Add CLI flags for different operations (index, search, serve)

	// Read the corpus directory
	files, err := os.ReadDir(PAGES_DIR)
	error_check(err)

	db, err := openDB()
	error_check(err)
	defer db.Close()

	hash := sha256.New()

	// Index each file
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", PAGES_DIR, file.Name())
		// get the metadata file
		// metadataFilePath := fmt.Sprintf("%s/%s", METADATA_DIR, file.Name())
		// metadataFile, err := os.ReadFile(metadataFilePath)
		// error_check(err)
		// var metadata DocMeta
		// err = json.Unmarshal(metadataFile, &metadata)
		// error_check(err)
		// fmt.Println("Metadata: ", metadata)

		if strings.HasSuffix(filePath, ".html") {
			hash.Write([]byte(filePath))
			docId := hash.Sum(nil)
			docMeta := index_file(filePath, file.Name(), docId, postings)
			err := saveDocMeta(db, docId, docMeta)

			if err != nil {
				fmt.Println("Error saving docmeta: ", err)
			} else {
				fmt.Println("Saved docmeta: ", docId)
			}
		} else {
			fmt.Println("Skipping: ", filePath)
		}
	}

	fmt.Printf("Total terms indexed: %d\n", len(postings))

	err = savePostings(db, postings)

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

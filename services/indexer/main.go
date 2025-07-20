package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

const CORPUS_DIR = "./corpus"

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
	files, err := os.ReadDir(CORPUS_DIR)
	error_check(err)

	db, err := openDB()
	error_check(err)
	defer db.Close()

	hash := sha256.New()

	// Index each file
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", CORPUS_DIR, file.Name())
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

	// TODO: Add minimal HTTP handler for search API
	search("python language", 10, db)
}

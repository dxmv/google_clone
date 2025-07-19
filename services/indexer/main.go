package main

import (
	"fmt"
	"os"
	"strings"
)

const CORPUS_DIR = "./corpus"

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	// TODO: Add CLI flags for different operations (index, search, serve)

	// Read the corpus directory
	files, err := os.ReadDir(CORPUS_DIR)
	error_check(err)
	docId := int64(0)

	db, err := openDB()
	error_check(err)
	defer db.Close()

	// Index each file
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", CORPUS_DIR, file.Name())
		if strings.HasSuffix(filePath, ".html") {
			docMeta := index_file(filePath, file.Name(), docId)
			saveDocMeta(db, docId, docMeta)
			docId++
		} else {
			fmt.Println("Skipping: ", filePath)
		}
	}

	// TODO: Add minimal HTTP handler for search API
}

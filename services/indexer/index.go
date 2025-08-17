package main

import (
	"fmt"
	"strings"

	shared "github.com/dxmv/google_clone/shared"
	"golang.org/x/net/html"
)

// addToIndex merges a document's term frequency map into the global inverted index
func addToIndex(docID []byte, termFreq map[string]int, postings map[string][]shared.Posting) {
	for term, count := range termFreq {
		_, ok := postings[term]
		if ok {
			postings[term] = append(postings[term], shared.Posting{DocID: docID, Count: count})
		} else {
			posting := shared.Posting{DocID: docID, Count: count}
			postings[term] = []shared.Posting{posting}
		}
	}
}

// Index an html file using the modular functions
func index_file(htmlString string, id []byte, postings map[string][]shared.Posting) {

	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(htmlString))
	error_check(err)

	// Extract text from HTML
	text := parse_html(doc)
	fmt.Println("Parsed html...\nTokenizing...")

	// Tokenize the text
	termFreq := shared.Tokenize(text)

	// add it to the index
	addToIndex(id, termFreq, postings)
	fmt.Println("Added to index...")

}

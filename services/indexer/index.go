package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Posting represents a document's relevance for a term
type Posting struct {
	DocID []byte
	Count int
}

// addToIndex merges a document's term frequency map into the global inverted index
func addToIndex(docID []byte, termFreq map[string]int, postings map[string][]Posting) {
	for term, count := range termFreq {
		_, ok := postings[term]
		if ok {
			postings[term] = append(postings[term], Posting{DocID: docID, Count: count})
		} else {
			posting := Posting{DocID: docID, Count: count}
			postings[term] = []Posting{posting}
		}
	}
}

// Index an html file using the modular functions
func index_file(filePath string, id []byte, postings map[string][]Posting) {
	fmt.Println("Indexing: ", filePath)
	content, err := os.ReadFile(filePath)
	error_check(err)

	// Parse HTML content
	html_string := string(content)
	doc, err := html.Parse(strings.NewReader(html_string))
	error_check(err)

	// Extract text from HTML
	text := parse_html(doc)
	fmt.Println("Parsed html...\nTokenizing...")

	// Tokenize the text
	termFreq := tokenise(text)

	// add it to the index
	addToIndex(id, termFreq, postings)
	fmt.Println("Added to index...")

	return
}

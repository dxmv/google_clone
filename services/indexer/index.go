package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// DocID represents document metadata
type DocMeta struct {
	Title    string
	FilePath string
	Length   int
}

// Posting represents a document's relevance for a term
type Posting struct {
	DocID int64
	Count int
}

// Global inverted index structure
var postings map[string][]Posting

// addToIndex merges a document's term frequency map into the global inverted index
func addToIndex(docID int64, termFreq map[string]int) {
	if len(postings) == 0 {
		postings = make(map[string][]Posting)
	}
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

// getPostings retrieves postings for a given term
func getPostings(term string) []Posting {
	return postings[term]
}

// Index an html file using the modular functions
func index_file(filePath string, fileName string, id int64) DocMeta {
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
	docMeta := DocMeta{Title: fileName, FilePath: filePath, Length: len(text)}
	termFreq := tokenise(text)

	// add it to the index
	addToIndex(id, termFreq)

	return docMeta
}

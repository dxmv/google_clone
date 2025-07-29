package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var START_LINKS = []string{
	"https://en.wikipedia.org/wiki/Philosophy",
}

type Job struct {
	URL   string
	Depth int
}

type DocMetadata struct {
	URL   string
	Depth int
	Title string
	// Meta_Description string
	Hash  string
	Links []string
}

func main() {
	createDirectory(PAGES_DIR)
	createDirectory(METADATA_DIR)

	visited := make(map[string]bool)
	queue := Queue[Job]{}
	// Add all start links to the queue
	for _, link := range START_LINKS {
		queue.Enqueue(Job{URL: link, Depth: 0})
	}
	// Process the queue
	for !queue.IsEmpty() {
		job := queue.Dequeue()
		if visited[job.URL] {
			continue
		}
		if job.Depth > 1 {
			break
		}
		visited[job.URL] = true
		fmt.Println("Visiting", job.URL)

		docMetadata := DocMetadata{
			URL:   job.URL,
			Depth: job.Depth,
			Title: "",
			Hash:  "",
		}
		// get the html
		body, err := fetch(job.URL)
		if err != nil {
			fmt.Println("Error getting HTML from", job.URL, err)
			continue
		}

		// extract the links from the html
		links := extractLinks(body, &docMetadata)
		docMetadata.Links = links
		for _, link := range links {
			queue.Enqueue(Job{URL: link, Depth: job.Depth + 1})
		}
		// save the html
		hash := sha256.Sum256(body)
		err = saveHTML(hex.EncodeToString(hash[:]), body)
		if err != nil {
			fmt.Println("Error saving HTML", err)
		}
		// save metadata
		docMetadata.Hash = hex.EncodeToString(hash[:])
		err = saveMetadata(docMetadata)
		if err != nil {
			fmt.Println("Error saving metadata", err)
		}

		fmt.Println("Visited", len(visited))
		fmt.Println("--------------------------------\n\n")
	}

	// print the visited links
	for link := range visited {
		fmt.Println(link)
	}
}

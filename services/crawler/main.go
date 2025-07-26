package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX_DEPTH = 2

var STARTING_LINKS = []string{
	"https://ubuntu.com",
}

var visited = make(map[string]bool)

func main() {
	// create the directories if they don't exist
	err := createDirectories()
	if err != nil {
		fmt.Println("Error creating directories: ", err)
		return
	}
	// do a simple BFS from starting links
	queue := JobQueue{}
	for _, link := range STARTING_LINKS {
		queue.Enqueue(Job{URL: link, Depth: 0})
	}

	for !queue.IsEmpty() {
		job := queue.Dequeue()
		// only go to max depth
		if job.Depth > MAX_DEPTH {
			continue
		}
		visited[job.URL] = true
		fmt.Println("Visiting: ", job)
		// get the html from the url
		delay := 3*time.Second + time.Duration(rand.Intn(500))*time.Millisecond
		time.Sleep(delay)
		body, err := getHtmlFromURL(job.URL)
		if err != nil {
			fmt.Println("Error getting HTML from URL: ", err)
			continue
		}
		// get links from the page and enqueue them
		parseResult := parsePage(body, job.URL)
		// save the page to a file
		savePage(job.URL, job.Depth, parseResult, body)
		for _, link := range parseResult.Links {
			// dont enqueue if already visited
			if visited[link] {
				continue
			}
			queue.Enqueue(Job{URL: link, Depth: job.Depth + 1})
		}
		fmt.Println("Visited: ", len(visited))
		fmt.Println("--------------------------------\n\n\n")
	}
	fmt.Println("\n\n\n\n--------------------------------")
	// print all visited links
	for link := range visited {
		fmt.Println(link)
	}
}

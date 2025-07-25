package main

import "fmt"

const MAX_DEPTH = 1

var STARTING_LINKS = []string{
	"https://sydneychiropractorcbd.com.au",
}

var visited = make(map[string]bool)

func main() {
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
		body, err := getHtmlFromURL(job.URL)
		if err != nil {
			fmt.Println("Error getting HTML from URL: ", err)
			continue
		}
		// get links from the page and enqueue them
		links := getLinks(body)
		// save the page to a file
		savePage(job.URL, job.Depth, len(links), body)
		for _, link := range links {
			// dont enqueue if already visited
			if visited[link] {
				continue
			}
			queue.Enqueue(Job{URL: link, Depth: job.Depth + 1})
		}
	}
	fmt.Println("\n\n\n\n--------------------------------")
	fmt.Println("Visited: ", len(visited))
	// print all visited links
	for link := range visited {
		fmt.Println(link)
	}
}

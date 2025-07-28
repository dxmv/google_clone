package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var START_LINKS = []string{
	"https://en.wikipedia.org/wiki/Wikipedia:Vital_articles/Level/1",
}

const PAGES_DIR = "pages"

type Job struct {
	URL   string
	Depth int
}

// Returns the body of the page
func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Check if the response is successful
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func handleHref(href string) (string, error) {
	if href == "" {
		return "", fmt.Errorf("href is empty")
	}
	if strings.HasPrefix(href, "#") {
		return "", fmt.Errorf("href " + href + " is a fragment")
	}
	res := ""
	if strings.HasPrefix(href, "/") {
		res = "https://en.wikipedia.org" + href
	}
	return res, nil
}

// Returns the links on the page
func extractLinks(body []byte) []string {

	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href, err := handleHref(attr.Val)
					if err != nil {
						fmt.Println("Error handling href", attr.Val, err)
						continue
					}
					links = append(links, href)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing", err)
		return []string{}
	}
	traverse(doc)
	return links
}

func saveHTML(hash string, body []byte) error {
	path := PAGES_DIR + "/" + hash + ".html"
	return os.WriteFile(path, body, 0644)
}

func createDirectory(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return os.Mkdir(name, 0755)
	}
	fmt.Println("Directory", name, "already exists")
	return nil
}

func main() {
	createDirectory(PAGES_DIR)

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
		body, err := fetch(job.URL)
		if err != nil {
			fmt.Println("Error getting HTML from", job.URL, err)
			continue
		}
		links := extractLinks(body)
		for _, link := range links {
			queue.Enqueue(Job{URL: link, Depth: job.Depth + 1})
		}
		hash := sha256.Sum256(body)
		err = saveHTML(hex.EncodeToString(hash[:]), body)
		if err != nil {
			fmt.Println("Error saving HTML", err)
		}
		fmt.Println("Visited", len(visited))
		fmt.Println("--------------------------------\n\n")
	}

	// print the visited links
	for link := range visited {
		fmt.Println(link)
	}
}

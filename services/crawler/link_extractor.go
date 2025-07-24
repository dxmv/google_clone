package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func getLinks(url string) []string {
	resp, err := http.Get(url)
	// get the html of page
	if err != nil {
		fmt.Println("Error fetching URL: ", err)
		return []string{}
	}
	defer resp.Body.Close()
	// skip if not 200 or not html
	if resp.StatusCode != 200 {
		fmt.Println("Error fetching URL: ", resp.StatusCode)
		return []string{}
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		fmt.Println("Not a HTML page")
		return []string{}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return []string{}
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML: ", err)
		return []string{}
	}
	// walk the html and get all links
	var links []string
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

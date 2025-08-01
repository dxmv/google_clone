package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

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
	if strings.HasPrefix(href, "#") {
		msg := fmt.Sprintf("href %s is a fragment", href)
		return "", errors.New(msg)
	}
	res := ""
	if strings.HasPrefix(href, "/") {
		res = "https://en.wikipedia.org" + href
	}
	return res, nil
}

// Returns the links on the page
func extractLinks(body []byte, docMetadata *DocMetadata) []string {

	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href, err := handleHref(attr.Val)
					if err != nil {
						fmt.Println("Error handling href", attr.Val, err)
						break
					}
					if href == "" {
						break
					}
					links = append(links, href)
				}
			}
		}
		// handle title
		if n.Type == html.ElementNode && n.Data == "title" {
			fmt.Println("Title", n.FirstChild.Data)
			docMetadata.Title = n.FirstChild.Data
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

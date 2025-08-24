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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoogleClone-Crawler/1.0 (Educational Project)")
	resp, err := http.DefaultClient.Do(req)
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

func handleImageSrc(src string) (string, error) {
	// only accept images from wikipedia and skip the wikipedia logo
	if !strings.Contains(src, "upload.wikimedia.org") {
		return "", errors.New("image src is not from wikipedia or is the wikipedia logo")
	}
	if strings.HasPrefix(src, "//") {
		src = "https:" + src
	} else if strings.HasPrefix(src, "/") {
		src = "https://en.wikipedia.org" + src
	}
	return src, nil
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
						continue
					}
					if href == "" {
						continue
					}
					links = append(links, href)
				}
			}
		}
		// images - extract src attribute from img tags
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imageUrl, err := handleImageSrc(attr.Val)
					if err != nil {
						continue
					}
					docMetadata.Images = append(docMetadata.Images, imageUrl)
					break
				}
			}
		}
		// handle title
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
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

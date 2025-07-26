package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type ParseResult struct {
	Links           []string
	Title           string
	MetaDescription string
}

func parsePage(body []byte, url string) ParseResult {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML: ", err)
		return ParseResult{}
	}
	// walk the html and get all links
	var links []string
	var title string
	var metaDescription string
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		// extract links
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// skip mailto & empty links
					if strings.HasPrefix(attr.Val, "mailto:") || attr.Val == "" || attr.Val[0] == '#' {
						continue
					}
					// handle relative links
					if attr.Val[0] == '/' {
						attr.Val = url + attr.Val
					}
					// handle relative links
					links = append(links, attr.Val)
				}
			}
		}
		// extract meta description
		if n.Type == html.ElementNode && n.Data == "meta" {
			for i, attr := range n.Attr {
				if attr.Key == "name" && (attr.Val == "description" || attr.Val == "Description") {
					if len(n.Attr) > i+1 {
						metaDescription = n.Attr[i+1].Val
					}
				}
			}
		}
		// extract meta title
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
		}
		// walk the children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return ParseResult{
		Links:           links,
		Title:           title,
		MetaDescription: metaDescription,
	}
}

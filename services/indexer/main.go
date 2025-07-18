package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const CORPUS_DIR = "./corpus"

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Parse the html, basically DFS on the DOM tree
func parse_html(doc *html.Node) string {
	var b strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			data := strings.TrimSpace(n.Data)
			if data != "" {
				b.WriteString(data)
				b.WriteByte(' ')
			}
		case html.ElementNode:
			if n.Data == "script" || n.Data == "style" {
				return // skip non-visible content
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return b.String()
}

// Index an html file
func index_file(filePath string) {
	fmt.Println("Indexing: ", filePath)
	content, err := os.ReadFile(filePath)
	error_check(err)
	// Format the html
	html_string := string(content)
	doc, err := html.Parse(strings.NewReader(html_string))
	error_check(err)

	html_string = parse_html(doc)
	fmt.Println("Parsed html...\n Tokenizing...")

	var m map[string]int
	m = make(map[string]int)
	for _, word := range strings.Split(html_string, " ") {
		m[word]++
	}

	fmt.Println(m)

}

func main() {
	// Read the corpus directory
	files, err := os.ReadDir(CORPUS_DIR)
	error_check(err)
	// Index each file
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", CORPUS_DIR, file.Name())
		if strings.HasSuffix(filePath, ".html") && strings.Contains(filePath, "Earth") {
			index_file(filePath)
		} else {
			fmt.Println("Skipping: ", filePath)
		}
	}
}

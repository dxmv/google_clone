package main

import (
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// Common stop words to filter out during indexing
var COMMON_WORDS map[string]bool = map[string]bool{
	"the":  true,
	"and":  true,
	"of":   true,
	"to":   true,
	"in":   true,
	"on":   true,
	"at":   true,
	"for":  true,
	"by":   true,
	"with": true,
	"as":   true,
	"from": true,
	"is":   true,
	"are":  true,
	"was":  true,
	"were": true,
	"be":   true,
	"been": true,
}

// isSeparator returns true for runes that should *split* tokens.
// - Letters → false  (stay inside the token)
// - Digits  → true   (acts as a boundary, so numbers are skipped)
// - Everything else (punctuation, whitespace, symbols) → true
func isSeparator(r rune) bool {
	return !unicode.IsLetter(r)
}

// Parse the html, basically DFS on the DOM tree
func parse_html(doc *html.Node) string {
	var b strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			data := strings.TrimSpace(n.Data)
			if data != "" && html.UnescapeString(data) != "" {
				b.WriteString(data)
				b.WriteByte(' ')
			}
		case html.ElementNode:
			if n.Data == "script" || n.Data == "style" || n.Data == "head" {
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

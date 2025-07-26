package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Metadata struct {
	URL             string `json:"url"`
	Depth           int    `json:"depth"`
	LengthOfLinks   int    `json:"length_of_links"`
	Title           string `json:"title"`
	MetaDescription string `json:"meta_description"`
	ContentLength   int    `json:"content_length"`
}

func getHtmlFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	// get the html of page
	if err != nil {
		fmt.Println("Error fetching URL: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	// skip if not 200 or not html
	if resp.StatusCode != 200 {
		fmt.Println("Error fetching URL: ", resp.StatusCode)
		return nil, err
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		fmt.Println("Not a HTML page")
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return nil, err
	}
	return body, nil
}

func createDirectories() error {
	err := ensureDir("pages")
	if err != nil {
		return err
	}
	err = ensureDir("metadata")
	if err != nil {
		return err
	}
	return nil
}

func savePage(url string, depth int, parseResult ParseResult, body []byte) error {
	// hash the url
	hash := sha256.Sum256([]byte(url))
	hashString := hex.EncodeToString(hash[:])

	fmt.Println("Hash: ", hashString)
	// save the html to <hash>.html in pages directory
	err := saveBytesToFile("pages/"+hashString+".html", body)
	if err != nil {
		fmt.Println("Error saving HTML to file: ", err)
		return err
	}
	// save the metadata to <hash>.json in metadata directory
	// first we need to create the metadata
	metadata := Metadata{
		URL:             url,
		Depth:           depth,
		LengthOfLinks:   len(parseResult.Links),
		Title:           parseResult.Title,
		MetaDescription: parseResult.MetaDescription,
		ContentLength:   len(body),
	}
	json, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling metadata: ", err)
		return err
	}
	err = saveBytesToFile("metadata/"+hashString+".json", json)
	if err != nil {
		fmt.Println("Error saving metadata to file: ", err)
		return err
	}
	return nil
}

func saveBytesToFile(path string, bytes []byte) error {
	return os.WriteFile(path, bytes, 0644)
}

func ensureDir(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("Directory already exists")
		} else {
			fmt.Println("Error creating directory: ", err)
			return err
		}
	}
	return nil
}

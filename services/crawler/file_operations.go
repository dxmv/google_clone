package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const PAGES_DIR = "pages"
const METADATA_DIR = "metadata"

func createDirectory(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return os.Mkdir(name, 0755)
	}
	fmt.Println("Directory", name, "already exists")
	return nil
}

func saveHTML(hash string, body []byte) error {
	path := PAGES_DIR + "/" + hash + ".html"
	return os.WriteFile(path, body, 0644)
}

func saveMetadata(docMetadata DocMetadata) error {
	path := METADATA_DIR + "/" + docMetadata.Hash + ".json"
	json, err := json.Marshal(docMetadata)
	if err != nil {
		return err
	}
	return os.WriteFile(path, json, 0644)
}

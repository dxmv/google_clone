package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage interface {
	SaveHTML(hash string, body []byte) error
	SaveMetadata(docMetadata DocMetadata) error
	CreateDirectory(name string) error
}

type LocalStorage struct {
	pagesDir    string
	metadataDir string
}

const PAGES_DIR = "pages1"
const METADATA_DIR = "metadata1"

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		pagesDir:    PAGES_DIR,
		metadataDir: METADATA_DIR,
	}
}

func (s *LocalStorage) CreateDirectory(name string) error {
	// check if the directory exists
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err := os.Mkdir(name, 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("Directory", name, "already exists")
	return nil
}

func (s *LocalStorage) SaveHTML(hash string, body []byte) error {
	path := PAGES_DIR + "/" + hash + ".html"
	return os.WriteFile(path, body, 0644)
}

func (s *LocalStorage) SaveMetadata(docMetadata DocMetadata) error {
	path := METADATA_DIR + "/" + docMetadata.Hash + ".json"
	json, err := json.Marshal(docMetadata)
	if err != nil {
		return err
	}
	return os.WriteFile(path, json, 0644)
}

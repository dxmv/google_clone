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

func NewLocalStorage(pagesDir string, metadataDir string) *LocalStorage {
	return &LocalStorage{
		pagesDir:    pagesDir,
		metadataDir: metadataDir,
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
	path := s.pagesDir + "/" + hash + ".html"
	return os.WriteFile(path, body, 0644)
}

func (s *LocalStorage) SaveMetadata(docMetadata DocMetadata) error {
	path := s.metadataDir + "/" + docMetadata.Hash + ".json"
	json, err := json.Marshal(docMetadata)
	if err != nil {
		return err
	}
	return os.WriteFile(path, json, 0644)
}

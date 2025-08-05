package main

import "runtime"

const PAGES_DIR = "pages"
const METADATA_DIR = "metadata"

type Config struct {
	StartLinks  []string
	MaxDepth    int
	JobsBuffer  int
	MaxRounds   int
	NumWorkers  int
	PagesDir    string
	MetadataDir string
}

func NewConfig() *Config {
	return &Config{
		StartLinks: []string{
			"https://en.wikipedia.org/wiki/Philosophy",
			"https://en.wikipedia.org/wiki/Mathematics",
		},
		MaxDepth:    1,
		JobsBuffer:  1000,
		MaxRounds:   10,
		NumWorkers:  runtime.NumCPU(),
		PagesDir:    PAGES_DIR,
		MetadataDir: METADATA_DIR,
	}
}

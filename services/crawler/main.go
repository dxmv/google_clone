package main

import (
	"fmt"
)

func main() {
	config := NewConfig()
	crawler := NewCrawler(NewLocalStorage(config.PagesDir, config.MetadataDir), config)
	err := crawler.Start()
	if err != nil {
		fmt.Println("Error starting crawler", err)
		return
	}
}

package main

import (
	"context"
	"fmt"
)

func main() {
	config := NewConfig()
	storage := NewMinioMongoStorage(config.MongoUri, config.MinioClient, context.Background())
	crawler := NewCrawler(storage, config)
	err := crawler.Start()
	if err != nil {
		fmt.Println("Error starting crawler", err)
		return
	}
}

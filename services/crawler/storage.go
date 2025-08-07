package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	SaveHTML(hash string, body []byte) error
	SaveMetadata(docMetadata DocMetadata) error
	CreateDirectory(name string) error
}

// minio and mongodb storage
type MinioMongoStorage struct {
	mongoConnection *mongo.Client
	pagesDir        string
	metadataDir     string
}

func newMongoConnection(uri string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMinioMongoStorage(mongoUri string, ctx context.Context) *MinioMongoStorage {
	mongoConnection, err := newMongoConnection(mongoUri, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return &MinioMongoStorage{
		mongoConnection: mongoConnection,
		pagesDir:        "pages",
		metadataDir:     "metadata",
	}
}

func (s *MinioMongoStorage) CreateDirectory(name string) error {
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

func (s *MinioMongoStorage) SaveHTML(hash string, body []byte) error {
	path := s.pagesDir + "/" + hash + ".html"
	return os.WriteFile(path, body, 0644)
}

func (s *MinioMongoStorage) SaveMetadata(docMetadata DocMetadata) error {
	coll := s.mongoConnection.Database("crawler").Collection("metadata")
	_, err := coll.InsertOne(context.Background(), docMetadata)
	if err != nil {
		return err
	}
	return nil
}

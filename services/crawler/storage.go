package main

import (
	"bytes"
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	SaveHTML(hash string, body []byte) error
	SaveMetadata(docMetadata DocMetadata) error
	CreateMetadataDirectory(name string) error
	CreateHTMLDirectory(name string) error
}

// minio and mongodb storage
type MinioMongoStorage struct {
	mongoConnection *mongo.Client
	minioClient     *minio.Client
}

func NewMinioMongoStorage(mongoUri string, minioClient *minio.Client, ctx context.Context) *MinioMongoStorage {
	mongoConnection, err := newMongoConnection(mongoUri, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return &MinioMongoStorage{
		mongoConnection: mongoConnection,
		minioClient:     minioClient,
	}
}

func (s *MinioMongoStorage) CreateMetadataDirectory(name string) error {
	// create a collection in mongodb
	coll := s.mongoConnection.Database("crawler").Collection(name)
	_, err := coll.InsertOne(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	return nil
}

// ensure that the bucket exists
func (s *MinioMongoStorage) CreateHTMLDirectory(name string) error {
	exists, err := s.minioClient.BucketExists(context.Background(), name)
	if err != nil {
		return err
	}
	if !exists {
		err = s.minioClient.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MinioMongoStorage) SaveHTML(hash string, body []byte) error {
	objectName := hash + ".html"
	contentType := "text/html"

	// upload the html to minio
	_, err := s.minioClient.PutObject(context.Background(), PAGES_DIR, objectName, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (s *MinioMongoStorage) SaveMetadata(docMetadata DocMetadata) error {
	coll := s.mongoConnection.Database("crawler").Collection("metadata")
	_, err := coll.InsertOne(context.Background(), docMetadata)
	if err != nil {
		return err
	}
	return nil
}

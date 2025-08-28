package main

import (
	"bytes"
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	SaveHTML(hash string, body []byte) error
	SaveMetadata(docMetadata DocMetadata) error
	CreateMetadataDirectory(name string) error
	CreateHTMLDirectory(name string) error
	FlushMetadata() error
}

// minio and mongodb storage
type MinioMongoStorage struct {
	mongoConnection *mongo.Client
	minioClient     *minio.Client
	metadataQueue   []interface{}
	maxMetadataJobs int
}

func NewMinioMongoStorage(mongoUri string, minioClient *minio.Client, ctx context.Context) *MinioMongoStorage {
	mongoConnection, err := newMongoConnection(mongoUri, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return &MinioMongoStorage{
		mongoConnection: mongoConnection,
		minioClient:     minioClient,
		metadataQueue:   make([]interface{}, 0),
		maxMetadataJobs: 300,
	}
}

func (s *MinioMongoStorage) CreateMetadataDirectory(name string) error {
	// create a collection in mongodb
	s.mongoConnection.Database("crawler").Collection(name)
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
		return err
	}
	return nil
}

func (s *MinioMongoStorage) SaveMetadata(docMetadata DocMetadata) error {
	s.metadataQueue = append(s.metadataQueue, docMetadata)
	if len(s.metadataQueue) >= s.maxMetadataJobs {
		return s.saveBatchMetadata()
	}
	return nil
}

func (s *MinioMongoStorage) saveBatchMetadata() error {
	if len(s.metadataQueue) == 0 {
		return nil
	}

	coll := s.mongoConnection.Database("crawler").Collection("metadata")
	models := make([]mongo.WriteModel, 0, len(s.metadataQueue))
	for _, it := range s.metadataQueue {
		m := it.(DocMetadata) // adjust if you store as interface{}
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"hash": m.Hash}). // or "url": m.URL
			SetUpdate(bson.M{"$set": m}).      // full doc update
			SetUpsert(true),
		)
	}
	_, err := coll.BulkWrite(context.Background(), models, options.BulkWrite().SetOrdered(false))
	s.metadataQueue = s.metadataQueue[:0]
	return err
}

func (s *MinioMongoStorage) FlushMetadata() error {
	if len(s.metadataQueue) == 0 {
		return nil
	}
	return s.saveBatchMetadata()
}

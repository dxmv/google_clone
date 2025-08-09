package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	MongoUri    string
	MinioClient *minio.Client
}

func NewConfig() *Config {
	monogUri := os.Getenv("MONGO_CONNECTION")
	minioClient, err := newMinioConnection()
	if err != nil {
		log.Println("Error creating minio client", err)
		return nil
	}
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
		MongoUri:    monogUri,
		MinioClient: minioClient,
	}
}

func newMongoConnection(uri string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newMinioConnection() (*minio.Client, error) {
	// read env variables
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	// create minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

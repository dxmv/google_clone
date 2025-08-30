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

const PAGES_DIR = "crawler-pages"
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
	}

	return &Config{
		StartLinks: []string{
			"https://en.wikipedia.org/wiki/Philosophy",
			"https://en.wikipedia.org/wiki/Mathematics",
			"https://en.wikipedia.org/wiki/Computer_science",
			"https://en.wikipedia.org/wiki/Economics",
			"https://en.wikipedia.org/wiki/Business",
			"https://en.wikipedia.org/wiki/Finance",
			"https://en.wikipedia.org/wiki/Astronomy",
			"https://en.wikipedia.org/wiki/Biology",
			"https://en.wikipedia.org/wiki/Chemistry",
			"https://en.wikipedia.org/wiki/Literature",
			"https://en.wikipedia.org/wiki/Physics",
			"https://en.wikipedia.org/wiki/Psychology",
			"https://en.wikipedia.org/wiki/Stock_market",
		},
		MaxDepth:    1,
		JobsBuffer:  10000,
		MaxRounds:   1000,
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

func newR2Client() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT") // e.g. 0aef...r2.cloudflarestorage.com
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       true,                  // R2 is HTTPS-only
		Region:       "auto",                // R2 expects "auto"
		BucketLookup: minio.BucketLookupDNS, // virtual-hosted style works best
	})
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

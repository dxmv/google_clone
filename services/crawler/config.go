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
	minioClient, err := newR2Client()
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

			"https://en.wikipedia.org/wiki/History",
			"https://en.wikipedia.org/wiki/Art",
			"https://en.wikipedia.org/wiki/Music",
			"https://en.wikipedia.org/wiki/Religion",
			"https://en.wikipedia.org/wiki/Law",

			"https://en.wikipedia.org/wiki/Medicine",
			"https://en.wikipedia.org/wiki/Engineering",
			"https://en.wikipedia.org/wiki/Geography",
			"https://en.wikipedia.org/wiki/Computer_engineering",
			"https://en.wikipedia.org/wiki/Artificial_intelligence",
			"https://en.wikipedia.org/wiki/Statistics",

			"https://en.wikipedia.org/wiki/Film",
			"https://en.wikipedia.org/wiki/Theatre",
			"https://en.wikipedia.org/wiki/Architecture",
			"https://en.wikipedia.org/wiki/Photography",

			"https://en.wikipedia.org/wiki/Internet",
			"https://en.wikipedia.org/wiki/Information_technology",

			"https://en.wikipedia.org/wiki/Space_exploration",
			"https://en.wikipedia.org/wiki/Mathematical_model",
			"https://en.wikipedia.org/wiki/Entrepreneurship",
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

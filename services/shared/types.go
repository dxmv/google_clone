package shared

import (
	"context"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

// DocMetadata represents the metadata of a document
type DocMetadata struct {
	URL           string
	Depth         int
	Title         string
	Hash          string
	Images        []string
	ContentLength int
	CrawledAt     time.Time
}

// Posting represents a document's relevance for a term
type Posting struct {
	DocID     []byte
	Count     int
	Positions []int
}

// Stats represents the statistics of the index
type Stats struct {
	AvgDocLength float64
	TotalDocs    int
}

type Corpus interface {
	GetHTML(ctx context.Context, hash string) ([]byte, error)
	ListMetadata(ctx context.Context) ([]DocMetadata, error)
	GetMetadata(ctx context.Context, docID string) (DocMetadata, error)
	GetBatchMetadata(ctx context.Context, docIDs []string) ([]DocMetadata, error)
}

type MinoMongoCorpus struct {
	minioClient    *minio.Client
	mongoClient    *mongo.Client
	bucketName     string
	collectionName string
	databaseName   string
}

// storage is a wrapper around the badger db and the corpus
type Storage struct {
	DB     *badger.DB
	Corpus Corpus
}

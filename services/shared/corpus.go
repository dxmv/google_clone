package shared

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMinoMongoCorpus() *MinoMongoCorpus {
	minio, err := newMinioConnection()
	if err != nil {
		log.Panic("Could not connect to minio", err)
	}
	mongo, err := newMongoConnection(context.Background())
	if err != nil {
		log.Panic("Could not connect to mongo", err)
	}
	return &MinoMongoCorpus{
		minioClient:    minio,
		mongoClient:    mongo,
		bucketName:     "pages",
		collectionName: "metadata",
		databaseName:   "crawler",
	}
}

// mongo connection
func newMongoConnection(ctx context.Context) (*mongo.Client, error) {
	mongoUri := os.Getenv("MONGO_CONNECTION")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// minio connection
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

// get html from minio
func (c *MinoMongoCorpus) GetHTML(ctx context.Context, hash string) ([]byte, error) {
	data, err := c.minioClient.GetObject(ctx, c.bucketName, hash, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return io.ReadAll(data)
}

// list metadata from mongo
func (c *MinoMongoCorpus) ListMetadata(ctx context.Context) ([]DocMetadata, error) {
	coll := c.mongoClient.Database(c.databaseName).Collection(c.collectionName)
	// retrives all documents from the collection
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// get all documents
	var docs []DocMetadata
	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// return a single metadata doc
func (c *MinoMongoCorpus) GetMetadata(ctx context.Context, docID string) (DocMetadata, error) {
	coll := c.mongoClient.Database(c.databaseName).Collection(c.collectionName)
	var doc DocMetadata
	log.Println("docID", docID)
	err := coll.FindOne(ctx, bson.D{{"hash", docID}}).Decode(&doc)
	if err != nil {
		return DocMetadata{}, err
	}
	return doc, nil
}

func (c *MinoMongoCorpus) GetBatchMetadata(ctx context.Context, docIDs []string) ([]DocMetadata, error) {
	coll := c.mongoClient.Database(c.databaseName).Collection(c.collectionName)
	cursor, err := coll.Find(ctx, bson.D{{"hash", bson.D{{"$in", docIDs}}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []DocMetadata
	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

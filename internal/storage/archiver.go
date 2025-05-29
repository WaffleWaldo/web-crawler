package storage

import (
	"context"
	"time"

	"web-crawler/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WebPage represents a crawled web page
type WebPage struct {
	URL         string    `bson:"url"`
	Title       string    `bson:"title"`
	Content     string    `bson:"content"`
	Links       []string  `bson:"links"`
	CrawledAt   time.Time `bson:"crawled_at"`
	StatusCode  int       `bson:"status_code"`
	ContentType string    `bson:"content_type"`
}

// Archiver defines the interface for storing crawled pages
type Archiver interface {
	Store(ctx context.Context, page *WebPage) error
	Close(ctx context.Context) error
}

// MongoArchiver implements the Archiver interface using MongoDB
type MongoArchiver struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoArchiver creates a new MongoDB archiver
func NewMongoArchiver(uri string, cfg config.MongoDBConfig) (*MongoArchiver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(cfg.MaxPoolSize)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(cfg.Database).Collection(cfg.Collection)
	return &MongoArchiver{
		client:     client,
		collection: collection,
	}, nil
}

// Store saves a webpage to MongoDB
func (m *MongoArchiver) Store(ctx context.Context, page *WebPage) error {
	_, err := m.collection.InsertOne(ctx, page)
	return err
}

// Close closes the MongoDB connection
func (m *MongoArchiver) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

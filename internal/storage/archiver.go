package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"web-crawler/internal/config"
	"web-crawler/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	logger.Info("Initializing MongoDB connection...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	// Configure MongoDB client options
	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxIdleTime).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetServerSelectionTimeout(cfg.Timeout).
		SetConnectTimeout(cfg.Timeout).
		SetSocketTimeout(cfg.Timeout).
		SetTLSConfig(tlsConfig).
		SetDirect(false).
		SetCompressors([]string{"snappy"}).
		SetReadPreference(readpref.Primary()).
		SetHeartbeatInterval(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Error("Failed to create MongoDB client: %v", err)
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Ping the database to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		logger.Error("Failed to ping MongoDB: %v", err)
		// Close the client if ping fails
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()
		_ = client.Disconnect(closeCtx)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Success("Successfully connected to MongoDB")

	// Get collection and ensure index
	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	// Create unique index on URL field if it doesn't exist
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		// If error is not because index already exists, return error
		if !mongo.IsDuplicateKeyError(err) {
			logger.Error("Failed to create index: %v", err)
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	}

	logger.Info("Using database: %s, collection: %s", cfg.Database, cfg.Collection)
	return &MongoArchiver{
		client:     client,
		collection: collection,
	}, nil
}

// Store saves a webpage to MongoDB using upsert
func (m *MongoArchiver) Store(ctx context.Context, page *WebPage) error {
	filter := bson.M{"url": page.URL}
	update := bson.M{
		"$set": bson.M{
			"title":        page.Title,
			"content":      page.Content,
			"links":        page.Links,
			"crawled_at":   page.CrawledAt,
			"status_code":  page.StatusCode,
			"content_type": page.ContentType,
		},
	}
	opts := options.Update().SetUpsert(true)

	result, err := m.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Failed to store/update webpage %s: %v", page.URL, err)
		return fmt.Errorf("failed to store/update webpage: %w", err)
	}

	// Log whether this was an insert or update
	if result.UpsertedCount > 0 {
		logger.StorageStatus(page.URL, false) // New document
	} else {
		logger.StorageStatus(page.URL, true) // Updated document
	}

	return nil
}

// Close closes the MongoDB connection
func (m *MongoArchiver) Close(ctx context.Context) error {
	logger.Info("Closing MongoDB connection...")
	if err := m.client.Disconnect(ctx); err != nil {
		logger.Error("Failed to disconnect from MongoDB: %v", err)
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	logger.Success("MongoDB connection closed")
	return nil
}

package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DefaultTimeout is the default timeout for database operations
const DefaultTimeout = 10 * time.Second

// WithTimeout creates a new context with the DefaultTimeout
func WithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

var globalDB *Database

func GetDatabase() (*Database, error) {
	if globalDB != nil {
		return globalDB, nil
	}

	connectionString, ok := os.LookupEnv("DB_URL")
	if !ok {
		panic("Environment variable DB_URL is not set")
	}

	ctx, cancel := WithTimeout()
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	globalDB = &Database{
		client: client,
		db:     client.Database("tranquil_pages"),
	}
	return globalDB, nil
}

func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

func (d *Database) Close() error {
	if d.client != nil {
		ctx, cancel := WithTimeout()
		defer cancel()
		if err := d.client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
			return err
		}
	}
	return nil
}

type TestDatabase struct {
	*Database
}

func NewTestDatabase() (*TestDatabase, error) {
	connectionString, ok := os.LookupEnv("DB_URL")
	if !ok {
		connectionString = "mongodb://localhost:27017" // Default localhost
	}

	ctx, cancel := WithTimeout()
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	testDBName := "tranquil_pages_test_" + uuid.New().String()
	testDB := client.Database(testDBName)

	err = testDB.CreateCollection(ctx, "books")
	if err != nil {
		return nil, err
	}

	return &TestDatabase{
		Database: &Database{
			client: client,
			db:     testDB,
		},
	}, nil
}

func (td *TestDatabase) Close() error {
	ctx, cancel := WithTimeout()
	defer cancel()

	if err := td.db.Drop(ctx); err != nil {
		log.Printf("Error dropping test database: %v", err)
	}

	if err := td.client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
		return err
	}

	return nil
}

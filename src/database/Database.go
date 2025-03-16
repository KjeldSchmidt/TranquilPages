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

var client *mongo.Client
var db *mongo.Database

func GetDbHandler() (*mongo.Database, error) {
	if db != nil {
		return db, nil
	}

	connectionString, ok := os.LookupEnv("DB_URL")
	if !ok {
		panic("Environment variable DB_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	db = client.Database("tranquil_pages")
	return db, nil
}

func CloseConnection() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

type TestDatabase struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewTestDatabase() (*TestDatabase, error) {
	connectionString, ok := os.LookupEnv("DB_URL")
	if !ok {
		connectionString = "mongodb://localhost:27017" // Default localhost
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		client: client,
		db:     testDB,
	}, nil
}

func (td *TestDatabase) GetDatabase() *mongo.Database {
	return td.db
}

func (td *TestDatabase) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

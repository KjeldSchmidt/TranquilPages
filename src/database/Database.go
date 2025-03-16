package database

import (
	"context"
	"log"
	"os"
	"time"

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

	db = client.Database("betterreads")
	return db, nil
}

func GetTestDatabase() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get the test database
	testDB := client.Database("tranquil_pages_test")

	// Drop the test database to ensure a clean state
	err = testDB.Drop(ctx)
	if err != nil {
		return nil, err
	}

	// Create the test database again
	testDB = client.Database("tranquil_pages_test")

	// Create the books collection
	err = testDB.CreateCollection(ctx, "books")
	if err != nil {
		return nil, err
	}

	return testDB, nil
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

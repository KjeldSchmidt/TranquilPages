package repository

import (
	"context"
	"errors"
	"log"
	"time"
	"tranquil-pages/src/database"
	appErrors "tranquil-pages/src/errors"
	"tranquil-pages/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	FindAll(ctx context.Context) ([]models.Book, error)
	FindById(ctx context.Context, id string) (*models.Book, error)
	Delete(ctx context.Context, id string) error
}

type MongoBookRepository struct {
	db *database.Database
}

func NewBookRepository(db *database.Database) BookRepository {
	return &MongoBookRepository{db: db}
}

func (r *MongoBookRepository) Create(ctx context.Context, book *models.Book) error {
	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := r.db.GetCollection("books").InsertOne(ctx, book)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return appErrors.ErrDuplicateBook
		}
		log.Printf("Database error in CreateBook: %v", err)
		return appErrors.ErrDatabase
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MongoBookRepository) FindAll(ctx context.Context) ([]models.Book, error) {
	cursor, err := r.db.GetCollection("books").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Database error in GetAllBooks: %v", err)
		return nil, appErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		log.Printf("Database error in GetAllBooks cursor.All: %v", err)
		return nil, appErrors.ErrDatabase
	}
	return books, nil
}

func (r *MongoBookRepository) FindById(ctx context.Context, id string) (*models.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, appErrors.ErrInvalidID
	}

	var book models.Book
	err = r.db.GetCollection("books").FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrNotFound
		}
		log.Printf("Database error in GetBookById: %v", err)
		return nil, appErrors.ErrDatabase
	}
	return &book, nil
}

func (r *MongoBookRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return appErrors.ErrInvalidID
	}

	_, err = r.db.GetCollection("books").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Database error in DeleteBook: %v", err)
		return appErrors.ErrDatabase
	}

	return nil
}

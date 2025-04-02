package repository

import (
	"errors"
	"log"
	"time"
	"tranquil-pages/database"
	appErrors "tranquil-pages/errors"
	"tranquil-pages/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	Create(book *models.Book) error
	FindAll() ([]models.Book, error)
	FindById(id string) (*models.Book, error)
	Delete(id string) error
	FindByUserID(userID string) ([]models.Book, error)
}

type MongoBookRepository struct {
	db *database.Database
}

func NewBookRepository(db *database.Database) BookRepository {
	return &MongoBookRepository{db: db}
}

func (r *MongoBookRepository) handleDBError(err error, operation string) error {
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return appErrors.ErrDuplicateBook
		}
		log.Printf("Database error in %s: %v", operation, err)
		return appErrors.ErrDatabase
	}
	return nil
}

func (r *MongoBookRepository) Create(book *models.Book) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := r.db.GetCollection("books").InsertOne(ctx, book)
	if err := r.handleDBError(err, "CreateBook"); err != nil {
		return err
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MongoBookRepository) FindAll() ([]models.Book, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	cursor, err := r.db.GetCollection("books").Find(ctx, bson.M{})
	if err := r.handleDBError(err, "GetAllBooks"); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err := r.handleDBError(cursor.All(ctx, &books), "GetAllBooks cursor.All"); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *MongoBookRepository) FindById(id string) (*models.Book, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

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
		return nil, r.handleDBError(err, "GetBookById")
	}
	return &book, nil
}

func (r *MongoBookRepository) Delete(id string) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return appErrors.ErrInvalidID
	}

	_, err = r.db.GetCollection("books").DeleteOne(ctx, bson.M{"_id": objectID})
	return r.handleDBError(err, "DeleteBook")
}

func (r *MongoBookRepository) FindByUserID(userID string) ([]models.Book, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	cursor, err := r.db.GetCollection("books").Find(ctx, bson.M{"user_id": userID})
	if err := r.handleDBError(err, "FindByUserID"); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err := r.handleDBError(cursor.All(ctx, &books), "FindByUserID cursor.All"); err != nil {
		return nil, err
	}
	return books, nil
}

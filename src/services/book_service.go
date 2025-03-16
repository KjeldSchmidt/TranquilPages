package services

import (
	"context"
	"errors"
	"time"
	appErrors "tranquil-pages/src/errors"
	"tranquil-pages/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookService struct {
	db *mongo.Database
}

func NewBookService(db *mongo.Database) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := s.db.Collection("books").InsertOne(ctx, book)
	if err != nil {
		return err
	}

	// Update the book with the new ID
	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

func (s *BookService) GetBookById(id string) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	var book models.Book
	err = s.db.Collection("books").FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &book, nil
}

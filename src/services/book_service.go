package services

import (
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

type BookService struct {
	db *database.Database
}

func NewBookService(db *database.Database) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *models.Book) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	if book.Rating < 0 || book.Rating > 5 {
		return appErrors.ErrInvalidRating
	}

	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := s.db.GetCollection("books").InsertOne(ctx, book)
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

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	cursor, err := s.db.GetCollection("books").Find(ctx, bson.M{})
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

func (s *BookService) GetBookById(id string) (*models.Book, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, appErrors.ErrInvalidID
	}

	var book models.Book
	err = s.db.GetCollection("books").FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrNotFound
		}
		log.Printf("Database error in GetBookById: %v", err)
		return nil, appErrors.ErrDatabase
	}
	return &book, nil
}

func (s *BookService) DeleteBook(id string) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return appErrors.ErrInvalidID
	}

	_, err = s.db.GetCollection("books").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Database error in DeleteBook: %v", err)
		return appErrors.ErrDatabase
	}

	return nil
}

package services

import (
	appErrors "betterreads/src/errors"
	"betterreads/src/models"
	"errors"
	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *models.Book) error {
	return s.db.Create(book).Error
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := s.db.Find(&books).Error
	return books, err
}

func (s *BookService) GetBookById(id string) (*models.Book, error) {
	var book *models.Book
	err := s.db.First(&book, id).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, appErrors.ErrNotFound
	case err != nil:
		return nil, err
	default:
		return book, err
	}
}

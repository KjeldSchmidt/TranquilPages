package services

import (
	"betterreads/src/models"
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

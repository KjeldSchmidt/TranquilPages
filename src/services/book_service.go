package services

import (
	"context"
	appErrors "tranquil-pages/src/errors"
	"tranquil-pages/src/models"
	"tranquil-pages/src/repository"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(book *models.Book) error {
	if book.Rating < 0 || book.Rating > 5 {
		return appErrors.ErrInvalidRating
	}

	return s.repo.Create(context.Background(), book)
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.FindAll(context.Background())
}

func (s *BookService) GetBookById(id string) (*models.Book, error) {
	return s.repo.FindById(context.Background(), id)
}

func (s *BookService) DeleteBook(id string) error {
	return s.repo.Delete(context.Background(), id)
}

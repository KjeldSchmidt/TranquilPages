package services

import (
	appErrors "tranquil-pages/errors"
	"tranquil-pages/models"
	"tranquil-pages/repository"
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

	return s.repo.Create(book)
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.FindAll()
}

func (s *BookService) GetBookById(id string) (*models.Book, error) {
	return s.repo.FindById(id)
}

func (s *BookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

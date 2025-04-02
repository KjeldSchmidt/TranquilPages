package services

import (
	"errors"
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

func (s *BookService) GetBooksByUserID(userID string) ([]models.Book, error) {
	return s.repo.FindByUserID(userID)
}

func (s *BookService) GetBook(id, userID string) (*models.Book, error) {
	book, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if book.UserID != userID {
		return nil, appErrors.ErrNotFound
	}
	return book, nil
}

func (s *BookService) DeleteBook(id, userID string) error {
	book, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil
		}
		return err
	}
	if book.UserID != userID {
		return appErrors.ErrNotFound
	}
	return s.repo.Delete(id)
}

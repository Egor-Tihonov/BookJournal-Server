package services

import (
	"book-journal/internal/models"
	"context"
	"net/http"
)

type BookService struct {
	repo   BookRepository
	client *http.Client
}

func NewBookService(repo BookRepository, client *http.Client) *BookService {
	return &BookService{repo: repo, client: client}
}

func (s *BookService) SearchBook(ctx context.Context, query models.Book) (*models.Book, error) {
	book, err := s.dbSearch(ctx, query)
	if err != nil {
		return nil, err // реальная ошибка БД
	}
	if book != nil {
		return book, nil // нашли в БД
	}

	return s.openLibrarySearch(ctx, query) // фоллбэк на внешний API
}

func (s *BookService) dbSearch(ctx context.Context, book models.Book) (*models.Book, error) {
	if book.ISBN != "" {
		return s.repo.FindByISBN(ctx, book.ISBN)
	}
	return s.repo.FindByName(ctx, book.Name)
}

func (s *BookService) openLibrarySearch(ctx context.Context, book models.Book) (*models.Book, error) {
	// https://openlibrary.org/api/books?bibkeys=ISBN:9781250319180&format=json&jscmd=data
	return nil, nil
}

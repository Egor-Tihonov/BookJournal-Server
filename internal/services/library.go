package services

import (
	"book-journal/internal/models"
	"context"
	"net/http"
)

// LibraryService owns the global book catalog (books table). It resolves a book
// by looking in our own catalog first and falling back to openlibrary. It knows
// nothing about users or their shelves.
type LibraryService struct {
	repo   LibraryRepository
	client *http.Client
}

func NewLibraryService(repo LibraryRepository, client *http.Client) *LibraryService {
	return &LibraryService{repo: repo, client: client}
}

// SearchBook выполняет поиск книги по ISBN или названию, сначала в базе данных, затем через внешний API OpenLibrary
func (s *LibraryService) SearchBook(ctx context.Context, query models.Book) (*models.Book, error) {
	book, err := s.dbSearch(ctx, query)
	if err != nil {
		return nil, err // реальная ошибка БД
	}
	if book != nil {
		return book, nil // нашли в БД
	}

	return s.openLibrarySearch(ctx, query) // фоллбэк на внешний API
}

// dbSearch выполняет поиск книги в базе данных по ISBN или названию
func (s *LibraryService) dbSearch(ctx context.Context, book models.Book) (*models.Book, error) {
	if book.ISBN != "" {
		return s.repo.FindByISBN(ctx, book.ISBN)
	}
	return s.repo.FindByName(ctx, book.Name)
}

// openLibrarySearch выполняет поиск книги через внешний API OpenLibrary
func (s *LibraryService) openLibrarySearch(ctx context.Context, book models.Book) (*models.Book, error) {
	// https://openlibrary.org/api/books?bibkeys=ISBN:9781250319180&format=json&jscmd=data
	return nil, nil
}

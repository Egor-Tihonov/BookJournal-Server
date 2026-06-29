package services

import (
	"context"

	"book-journal/internal/models"
)

// AddBookInput is what a user submits to put a book on their shelf.
type AddBookInput struct {
	ISBN       string
	Title      string
	Status     models.ReadingStatus
	WhyReading string
}

// LibraryService orchestrates a user's shelf (user_library). It does not search
// for books itself — it asks BookService to resolve the book, then records it
// on the user's shelf.
type LibraryService struct {
	books *BookService
	repo  LibraryRepository
}

func NewLibraryService(books *BookService, repo LibraryRepository) *LibraryService {
	return &LibraryService{books: books, repo: repo}
}

// AddBookToUserLibrary resolves the book (catalog → openlibrary) and adds it to
// the user's shelf. Returns ErrBookNotFound if the book exists nowhere.
func (s *LibraryService) AddBookToUserLibrary(ctx context.Context, userID int64, in AddBookInput) (*models.LibraryEntry, error) {
	book, err := s.books.SearchBook(ctx, models.Book{ISBN: in.ISBN, Name: in.Title})
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, ErrBookNotFound
	}

	status := in.Status
	if status == "" {
		status = models.StatusWantToRead
	}

	return s.repo.Add(ctx, userID, book.ID, status, in.WhyReading)
}

// GetLibrary returns every book on the user's shelf.
func (s *LibraryService) GetLibrary(ctx context.Context, userID int64) ([]models.LibraryEntry, error) {
	return s.repo.ListByUser(ctx, userID)
}

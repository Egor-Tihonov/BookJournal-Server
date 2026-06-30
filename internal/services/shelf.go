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

// ShelfService orchestrates a user's personal shelf (user_shelf). It does not
// search for books itself — it asks LibraryService to resolve the book in the
// global catalog, then records it on the user's shelf.
type ShelfService struct {
	library *LibraryService
	repo    ShelfRepository
}

func NewShelfService(library *LibraryService, repo ShelfRepository) *ShelfService {
	return &ShelfService{library: library, repo: repo}
}

// AddBookToShelf resolves the book (catalog → openlibrary) and adds it to the
// user's shelf. Returns ErrBookNotFound if the book exists nowhere.
func (s *ShelfService) AddBookToShelf(ctx context.Context, userID int64, in AddBookInput) (*models.ShelfEntry, error) {
	book, err := s.library.SearchBook(ctx, models.Book{ISBN: in.ISBN, Name: in.Title})
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

// GetShelf returns every book on the user's shelf.
func (s *ShelfService) GetShelf(ctx context.Context, userID int64) ([]models.ShelfEntry, error) {
	return s.repo.ListByUser(ctx, userID)
}

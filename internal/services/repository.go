package services

import (
	"context"

	"book-journal/internal/models"
)

// LibraryRepository is the global catalog storage (books table) that
// LibraryService needs. Implemented by internal/db.
type LibraryRepository interface {
	FindByISBN(ctx context.Context, isbn string) (*models.Book, error)
	FindByName(ctx context.Context, name string) (*models.Book, error)
	Save(ctx context.Context, book *models.Book) (*models.Book, error)
}

// ShelfRepository is the storage for personal shelves (user_shelf, quotes,
// notes) that ShelfService needs. Implemented by internal/db.
type ShelfRepository interface {
	Add(ctx context.Context, userID, bookID int64, status models.ReadingStatus, whyReading string) (*models.ShelfEntry, error)
	ListByUser(ctx context.Context, userID int64) ([]models.ShelfEntry, error)
}

package services

import (
	"context"

	"book-journal/internal/models"
)

// BookRepository is the catalog storage (books table) that BookService needs.
// Implemented by internal/db.
type BookRepository interface {
	FindByISBN(ctx context.Context, isbn string) (*models.Book, error)
	FindByName(ctx context.Context, name string) (*models.Book, error)
	Save(ctx context.Context, book *models.Book) (*models.Book, error)
}

// LibraryRepository is the storage for user shelves (user_library, quotes,
// notes) that LibraryService needs. Implemented by internal/db.
type LibraryRepository interface {
	Add(ctx context.Context, userID, bookID int64, status models.ReadingStatus, whyReading string) (*models.LibraryEntry, error)
	ListByUser(ctx context.Context, userID int64) ([]models.LibraryEntry, error)
}

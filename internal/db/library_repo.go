package db

import (
	"context"

	"book-journal/internal/models"
)

// LibraryRepo is the database-backed user shelf (user_library, quotes, notes).
// TODO: back this with a real *sql.DB (Postgres). Methods are stubs for now.
type LibraryRepo struct{}

func NewLibraryRepo() *LibraryRepo {
	return &LibraryRepo{}
}

func (r *LibraryRepo) Add(ctx context.Context, userID, bookID int64, status models.ReadingStatus, whyReading string) (*models.LibraryEntry, error) {
	return nil, ErrNotImplemented
}

func (r *LibraryRepo) ListByUser(ctx context.Context, userID int64) ([]models.LibraryEntry, error) {
	return nil, ErrNotImplemented
}

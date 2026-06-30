package db

import (
	"context"

	"book-journal/internal/models"
)

// LibraryRepo is the database-backed global book catalog (books table).
// TODO: back this with a real *sql.DB (Postgres). Methods are stubs for now.
type LibraryRepo struct{}

func NewLibraryRepo() *LibraryRepo {
	return &LibraryRepo{}
}

func (r *LibraryRepo) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	return nil, ErrNotImplemented
}

func (r *LibraryRepo) FindByName(ctx context.Context, name string) (*models.Book, error) {
	return nil, ErrNotImplemented
}

func (r *LibraryRepo) Save(ctx context.Context, book *models.Book) (*models.Book, error) {
	return nil, ErrNotImplemented
}

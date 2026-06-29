package db

import (
	"context"

	"book-journal/internal/models"
)

// BookRepo is the database-backed book catalog (books table).
// TODO: back this with a real *sql.DB (Postgres). Methods are stubs for now.
type BookRepo struct{}

func NewBookRepo() *BookRepo {
	return &BookRepo{}
}

func (r *BookRepo) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	return nil, ErrNotImplemented
}

func (r *BookRepo) FindByName(ctx context.Context, name string) (*models.Book, error) {
	return nil, ErrNotImplemented
}

func (r *BookRepo) Save(ctx context.Context, book *models.Book) (*models.Book, error) {
	return nil, ErrNotImplemented
}

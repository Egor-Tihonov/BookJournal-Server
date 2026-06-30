package db

import (
	"context"

	"book-journal/internal/models"
)

// ShelfRepo is the database-backed personal shelf (user_shelf, quotes, notes).
// TODO: back this with a real *sql.DB (Postgres). Methods are stubs for now.
type ShelfRepo struct{}

func NewShelfRepo() *ShelfRepo {
	return &ShelfRepo{}
}

func (r *ShelfRepo) Add(ctx context.Context, userID, bookID int64, status models.ReadingStatus, whyReading string) (*models.ShelfEntry, error) {
	return nil, ErrNotImplemented
}

func (r *ShelfRepo) ListByUser(ctx context.Context, userID int64) ([]models.ShelfEntry, error) {
	return nil, ErrNotImplemented
}

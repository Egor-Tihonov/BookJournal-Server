package services

import "errors"

// ErrBookNotFound is returned when a book cannot be found in the catalog or via
// openlibrary.
var ErrBookNotFound = errors.New("book not found")

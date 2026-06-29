package db

import "errors"

// ErrNotImplemented marks repository methods that still need a real
// database-backed implementation (Postgres).
var ErrNotImplemented = errors.New("db: not implemented")

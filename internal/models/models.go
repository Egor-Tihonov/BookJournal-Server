package models

import "time"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// PasswordHash is never serialized to JSON.
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// Book is an entry in the global catalog (books table). It knows nothing
// about users.
type Book struct {
	ID     int64  `json:"id"`
	ISBN   string `json:"isbn,omitempty"`
	Name   string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
	Source string `json:"source,omitempty"` // manual / openlibrary
}

type ReadingStatus string

const (
	StatusWantToRead ReadingStatus = "want_to_read"
	StatusReading    ReadingStatus = "reading"
	StatusRead       ReadingStatus = "read"
)

// ShelfEntry is one book on a user's shelf (user_shelf) together with its quotes
// and notes. It represents the user's relationship with a book (status,
// thoughts, dates) — not the book's catalog metadata, which lives in Book.
type ShelfEntry struct {
	ID         int64         `json:"id"`
	UserID     int64         `json:"user_id"`
	Book       Book          `json:"book"`
	Status     ReadingStatus `json:"status"`
	WhyReading string        `json:"why_reading,omitempty"`
	StartedAt  *time.Time    `json:"started_at,omitempty"`
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
	Quotes     []Quote       `json:"quotes,omitempty"`
	Notes      []Note        `json:"notes,omitempty"`
}

type Quote struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Page      int       `json:"page,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Note struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

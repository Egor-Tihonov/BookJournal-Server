package api

import (
	"errors"
	"net/http"

	"book-journal/internal/models"
	"book-journal/internal/services"

	"github.com/labstack/echo"
)

// Handlers holds the dependencies the HTTP layer needs.
type Handlers struct {
	shelf *services.ShelfService
}

func NewHandlers(shelf *services.ShelfService) *Handlers {
	return &Handlers{shelf: shelf}
}

// addBookRequest is the JSON body for POST /book.
type addBookRequest struct {
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	WhyReading string `json:"why_reading"`
}

func (h *Handlers) AddBook(c echo.Context) error {
	var req addBookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	entry, err := h.shelf.AddBookToShelf(c.Request().Context(), currentUserID(c), services.AddBookInput{
		ISBN:       req.ISBN,
		Title:      req.Title,
		Status:     models.ReadingStatus(req.Status),
		WhyReading: req.WhyReading,
	})
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "book not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not add book")
	}

	return c.JSON(http.StatusCreated, entry)
}

func (h *Handlers) GetShelf(c echo.Context) error {
	entries, err := h.shelf.GetShelf(c.Request().Context(), currentUserID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not load shelf")
	}
	return c.JSON(http.StatusOK, entries)
}

// currentUserID returns the authenticated user's id.
// TODO: read this from the auth middleware once it is implemented.
func currentUserID(c echo.Context) int64 {
	return 1
}

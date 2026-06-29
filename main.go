package main

import (
	"net/http"

	"book-journal/api"
	"book-journal/internal/db"
	"book-journal/internal/services"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Storage + external clients (swap stubs for real implementations).
	bookRepo := db.NewBookRepo()
	libraryRepo := db.NewLibraryRepo()

	// Services: BookService owns the catalog; LibraryService orchestrates shelves.
	bookService := services.NewBookService(bookRepo, http.DefaultClient)
	libraryService := services.NewLibraryService(bookService, libraryRepo)

	// HTTP layer.
	handlers := api.NewHandlers(libraryService)
	api.RegisterRoutes(e, handlers)

	e.Logger.Fatal(e.Start(":8080"))
}

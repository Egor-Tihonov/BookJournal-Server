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
	libraryRepo := db.NewLibraryRepo()
	shelfRepo := db.NewShelfRepo()

	// Services: LibraryService owns the global catalog; ShelfService orchestrates
	// each user's personal shelf.
	libraryService := services.NewLibraryService(libraryRepo, http.DefaultClient)
	shelfService := services.NewShelfService(libraryService, shelfRepo)

	// HTTP layer.
	handlers := api.NewHandlers(shelfService)
	api.RegisterRoutes(e, handlers)

	e.Logger.Fatal(e.Start(":8080"))
}

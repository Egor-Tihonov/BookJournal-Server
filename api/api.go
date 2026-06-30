package api

import "github.com/labstack/echo"

// RegisterRoutes wires every endpoint in one place so the whole API surface is
// visible at a glance.
func RegisterRoutes(e *echo.Echo, h *Handlers) {
	api := e.Group("/api/v1")

	// TODO: middleware — authentication, authorization, request validation.

	api.GET("/shelf", h.GetShelf)
	api.POST("/book", h.AddBook)
}

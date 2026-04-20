package router

import (
	"github.com/gofiber/fiber/v2"

	"rate-limited/internal/handler"
)

// SetupRoutes sets up all API endpoints for the application
func SetupRoutes(app *fiber.App, h *handler.RequestHandler) {

	// Basic health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Rate-Limited API Service is running",
		})
	})

	// Main APIs
	app.Post("/request", h.HandleRequest)
	app.Get("/stats", h.GetStats)
}

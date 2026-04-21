package router

import (
	"github.com/gofiber/fiber/v2"

	"rate-limited/internal/handler"
)

func SetupRoutes(app *fiber.App, h *handler.RequestHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Rate-Limited API Service is running",
		})
	})

	app.Get("/health", h.Health)
	app.Post("/request", h.HandleRequest)
	app.Get("/stats", h.GetStats)
	app.Get("/jobs/:id", h.GetJob)
	app.Get("/queue/stats", h.GetQueueStats)
}

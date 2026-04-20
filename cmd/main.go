package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"rate-limited/config"
	"rate-limited/internal/handler"
	"rate-limited/internal/router"
	"rate-limited/internal/service"
	"rate-limited/internal/storage"
)

func main() {
	// Load environment configuration.
	cfg := config.LoadConfig()

	// Initialize in-memory store.
	store := storage.NewMemoryStore()

	// Initialize rate limiter service.
	rateLimiterService := service.NewRateLimiterService(
		store,
		cfg.RateLimitMaxRequests,
		cfg.RateLimitWindowSec,
	)

	// Initialize handlers.
	requestHandler := handler.NewRequestHandler(rateLimiterService)

	// Create Fiber app.
	app := fiber.New()

	// Register routes.
	router.SetupRoutes(app, requestHandler)

	// Start server.
	log.Fatal(app.Listen(":" + cfg.Port))
	log.Println("Server starting on port:", cfg.Port)
}

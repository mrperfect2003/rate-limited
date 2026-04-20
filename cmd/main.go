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
	// load config from .env (or system env if .env not present)
	cfg := config.LoadConfig()

	// in-memory store to keep request data (no DB as per assignment)
	store := storage.NewMemoryStore()

	// rate limiter service handles core logic
	rateLimiterService := service.NewRateLimiterService(
		store,
		cfg.RateLimitMaxRequests,
		cfg.RateLimitWindowSec,
	)

	// handler layer (API layer)
	h := handler.NewRequestHandler(rateLimiterService)

	// create fiber app
	app := fiber.New()

	// setup all routes
	router.SetupRoutes(app, h)

	// small log before starting server (helps during debugging)
	log.Println("Server starting on port:", cfg.Port)

	// start server (this blocks)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

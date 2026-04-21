package main

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"rate-limited/config"
	"rate-limited/internal/handler"
	"rate-limited/internal/model"
	"rate-limited/internal/queue"
	"rate-limited/internal/router"
	"rate-limited/internal/service"
	"rate-limited/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	var store storage.RateLimitStore

	switch cfg.StorageType {
	case "redis":
		log.Println("using redis store")
		store = storage.NewRedisStore(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	default:
		log.Println("using memory store")
		store = storage.NewMemoryStore()
	}

	jq := queue.NewJobQueue(
		cfg.QueueSize,
		cfg.MaxRetries,
		cfg.WorkerCount,
		func(job *model.Job) error {
			allowed := store.AllowRequest(job.UserID, cfg.RateLimitMaxRequests, time.Duration(cfg.RateLimitWindowSec)*time.Second)
			if !allowed {
				return errors.New("rate limit still active, retrying")
			}
			return nil
		},
	)

	rateLimiterService := service.NewRateLimiterService(
		store,
		cfg.RateLimitMaxRequests,
		cfg.RateLimitWindowSec,
		cfg.QueueEnabled,
		jq,
	)

	h := handler.NewRequestHandler(rateLimiterService)

	app := fiber.New()
	router.SetupRoutes(app, h)

	log.Println("server starting on port:", cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

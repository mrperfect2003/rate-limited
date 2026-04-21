package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	RateLimitMaxRequests int
	RateLimitWindowSec   int

	StorageType string
	RedisAddr   string
	RedisPass   string
	RedisDB     int

	QueueEnabled bool
	QueueSize    int
	MaxRetries   int
	WorkerCount  int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	maxRequests := mustAtoi("RATE_LIMIT_MAX_REQUESTS", "5")
	windowSec := mustAtoi("RATE_LIMIT_WINDOW_SECONDS", "60")
	redisDB := mustAtoi("REDIS_DB", "0")
	queueSize := mustAtoi("QUEUE_SIZE", "100")
	maxRetries := mustAtoi("MAX_RETRIES", "3")
	workerCount := mustAtoi("WORKER_COUNT", "2")

	return &Config{
		Port:                 getEnv("PORT", "5000"),
		RateLimitMaxRequests: maxRequests,
		RateLimitWindowSec:   windowSec,

		StorageType: getEnv("STORAGE_TYPE", "memory"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:     redisDB,

		QueueEnabled: mustBool("QUEUE_ENABLED", "true"),
		QueueSize:    queueSize,
		MaxRetries:   maxRetries,
		WorkerCount:  workerCount,
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustAtoi(key, fallback string) int {
	v, err := strconv.Atoi(getEnv(key, fallback))
	if err != nil {
		log.Fatalf("%s should be a valid integer", key)
	}
	return v
}

func mustBool(key, fallback string) bool {
	v, err := strconv.ParseBool(getEnv(key, fallback))
	if err != nil {
		log.Fatalf("%s should be true or false", key)
	}
	return v
}

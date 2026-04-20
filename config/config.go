package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the environment-based settings for the app
type Config struct {
	Port                 string
	RateLimitMaxRequests int
	RateLimitWindowSec   int
}

// LoadConfig loads values from .env (if present) and system env
func LoadConfig() *Config {
	// Try loading .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	// Convert rate limit values from string → int
	maxRequests, err := strconv.Atoi(getEnv("RATE_LIMIT_MAX_REQUESTS", "5"))
	if err != nil {
		log.Fatal("RATE_LIMIT_MAX_REQUESTS should be a valid number")
	}

	windowSec, err := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW_SECONDS", "60"))
	if err != nil {
		log.Fatal("RATE_LIMIT_WINDOW_SECONDS should be a valid number")
	}

	return &Config{
		Port:                 getEnv("PORT", "5000"),
		RateLimitMaxRequests: maxRequests,
		RateLimitWindowSec:   windowSec,
	}
}

// getEnv reads env variable, returns fallback if not present
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

package service

import (
	"time"

	"rate-limited/internal/model"
	"rate-limited/internal/storage"
)

// RateLimiterService contains business logic for rate limiting and stats.
type RateLimiterService struct {
	store       *storage.MemoryStore
	maxRequests int
	window      time.Duration
}

// NewRateLimiterService creates a new RateLimiterService.
func NewRateLimiterService(store *storage.MemoryStore, maxRequests int, windowSeconds int) *RateLimiterService {
	return &RateLimiterService{
		store:       store,
		maxRequests: maxRequests,
		window:      time.Duration(windowSeconds) * time.Second,
	}
}

// ProcessRequest validates whether a user request is allowed.
func (s *RateLimiterService) ProcessRequest(userID string) bool {
	return s.store.AllowRequest(userID, s.maxRequests, s.window)
}

// GetStats returns paginated stats for all users.
func (s *RateLimiterService) GetStats(page, limit int) model.StatsResponse {
	allStats := s.store.GetStats()
	totalUsers := len(allStats)

	start := (page - 1) * limit
	end := start + limit

	if start > totalUsers {
		start = totalUsers
	}
	if end > totalUsers {
		end = totalUsers
	}

	paginatedStats := allStats[start:end]

	return model.StatsResponse{
		Page:       page,
		Limit:      limit,
		TotalUsers: totalUsers,
		Stats:      paginatedStats,
	}
}

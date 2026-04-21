package storage

import (
	"time"

	"rate-limited/internal/model"
)

type RateLimitStore interface {
	AllowRequest(userID string, maxRequests int, window time.Duration) bool
	IncrementRejected(userID string)
	IncrementQueued(userID string)
	GetStats() []model.UserStat
}

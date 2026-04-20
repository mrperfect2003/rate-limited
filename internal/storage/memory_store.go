package storage

import (
	"sort"
	"sync"
	"time"

	"rate-limited/internal/model"
)

// MemoryStore keeps all request tracking data in memory.
// It is protected by a mutex so concurrent requests are handled safely.
type MemoryStore struct {
	mu sync.Mutex

	// requestLog stores timestamps of accepted requests per user.
	requestLog map[string][]time.Time

	// totalAccepted stores total accepted request count per user.
	totalAccepted map[string]int
}

// NewMemoryStore initializes and returns a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requestLog:    make(map[string][]time.Time),
		totalAccepted: make(map[string]int),
	}
}

// AllowRequest checks whether a user can make a request under the configured rate limit.
// It also updates in-memory state if the request is allowed.
//
// maxRequests: allowed number of requests in the time window
// window: duration for rate limit window, e.g. 1 minute
func (s *MemoryStore) AllowRequest(userID string, maxRequests int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	timestamps := s.requestLog[userID]

	// Keep only timestamps that are still inside the active window.
	validTimestamps := make([]time.Time, 0, len(timestamps))
	for _, ts := range timestamps {
		if now.Sub(ts) < window {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// If limit already reached, reject request.
	if len(validTimestamps) >= maxRequests {
		s.requestLog[userID] = validTimestamps
		return false
	}

	// Accept request and store current timestamp.
	validTimestamps = append(validTimestamps, now)
	s.requestLog[userID] = validTimestamps
	s.totalAccepted[userID]++

	return true
}

// GetStats returns all per-user accepted request counts sorted by user_id.
// Sorting makes pagination stable and predictable.
func (s *MemoryStore) GetStats() []model.UserStat {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := make([]model.UserStat, 0, len(s.totalAccepted))
	for userID, count := range s.totalAccepted {
		stats = append(stats, model.UserStat{
			UserID:        userID,
			TotalRequests: count,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].UserID < stats[j].UserID
	})

	return stats
}

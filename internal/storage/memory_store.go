package storage

import (
	"sort"
	"sync"
	"time"

	"rate-limited/internal/model"
)

type MemoryStore struct {
	mu sync.Mutex

	requestLog map[string][]time.Time
	stats      map[string]*model.UserStat
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requestLog: make(map[string][]time.Time),
		stats:      make(map[string]*model.UserStat),
	}
}

func (s *MemoryStore) getOrCreateUserStat(userID string) *model.UserStat {
	stat, exists := s.stats[userID]
	if !exists {
		stat = &model.UserStat{UserID: userID}
		s.stats[userID] = stat
	}
	return stat
}

func (s *MemoryStore) AllowRequest(userID string, maxRequests int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	timestamps := s.requestLog[userID]

	valid := make([]time.Time, 0, len(timestamps))
	for _, ts := range timestamps {
		if now.Sub(ts) < window {
			valid = append(valid, ts)
		}
	}

	stat := s.getOrCreateUserStat(userID)

	if len(valid) >= maxRequests {
		s.requestLog[userID] = valid
		stat.TotalRejected++
		return false
	}

	valid = append(valid, now)
	s.requestLog[userID] = valid
	stat.TotalAccepted++
	return true
}

func (s *MemoryStore) IncrementRejected(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.getOrCreateUserStat(userID)
	stat.TotalRejected++
}

func (s *MemoryStore) IncrementQueued(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.getOrCreateUserStat(userID)
	stat.QueuedRequests++
}

func (s *MemoryStore) GetStats() []model.UserStat {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := make([]model.UserStat, 0, len(s.stats))
	for _, stat := range s.stats {
		stats = append(stats, *stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].UserID < stats[j].UserID
	})

	return stats
}

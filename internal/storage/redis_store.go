package storage

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"rate-limited/internal/model"
)

type RedisStore struct {
	client *redis.Client

	mu    sync.Mutex
	stats map[string]*model.UserStat
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStore{
		client: rdb,
		stats:  make(map[string]*model.UserStat),
	}
}

func (s *RedisStore) getOrCreateUserStat(userID string) *model.UserStat {
	stat, exists := s.stats[userID]
	if !exists {
		stat = &model.UserStat{UserID: userID}
		s.stats[userID] = stat
	}
	return stat
}

func (s *RedisStore) AllowRequest(userID string, maxRequests int, window time.Duration) bool {
	ctx := context.Background()
	key := "rate_limit:" + userID
	now := time.Now().UnixMilli()
	windowStart := time.Now().Add(-window).UnixMilli()

	luaScript := `
redis.call("ZREMRANGEBYSCORE", KEYS[1], 0, ARGV[1])
local current = redis.call("ZCARD", KEYS[1])
if current >= tonumber(ARGV[2]) then
	return 0
end
redis.call("ZADD", KEYS[1], ARGV[3], ARGV[3])
redis.call("PEXPIRE", KEYS[1], ARGV[4])
return 1
`

	result, err := s.client.Eval(
		ctx,
		luaScript,
		[]string{key},
		windowStart,
		maxRequests,
		now,
		window.Milliseconds(),
	).Int()

	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.getOrCreateUserStat(userID)

	if err != nil {
		log.Println("redis eval error:", err)
		stat.TotalRejected++
		return false
	}

	if result == 1 {
		stat.TotalAccepted++
		return true
	}

	stat.TotalRejected++
	return false
}

func (s *RedisStore) IncrementRejected(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.getOrCreateUserStat(userID)
	stat.TotalRejected++
}

func (s *RedisStore) IncrementQueued(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.getOrCreateUserStat(userID)
	stat.QueuedRequests++
}

func (s *RedisStore) GetStats() []model.UserStat {
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

func (s *RedisStore) SaveStatsSnapshot() error {
	ctx := context.Background()
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s.stats)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, "stats_snapshot", string(data), 0).Err()
}

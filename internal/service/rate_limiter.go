package service

import (
	"errors"
	"time"

	"rate-limited/internal/model"
	"rate-limited/internal/queue"
	"rate-limited/internal/storage"
)

type RateLimiterService struct {
	store        storage.RateLimitStore
	maxRequests  int
	window       time.Duration
	queueEnabled bool
	jobQueue     *queue.JobQueue
}

func NewRateLimiterService(
	store storage.RateLimitStore,
	maxRequests int,
	windowSeconds int,
	queueEnabled bool,
	jobQueue *queue.JobQueue,
) *RateLimiterService {
	return &RateLimiterService{
		store:        store,
		maxRequests:  maxRequests,
		window:       time.Duration(windowSeconds) * time.Second,
		queueEnabled: queueEnabled,
		jobQueue:     jobQueue,
	}
}

func (s *RateLimiterService) ProcessRequest(userID string) bool {
	return s.store.AllowRequest(userID, s.maxRequests, s.window)
}

func (s *RateLimiterService) HandleIncomingRequest(userID string, payload interface{}) (bool, *model.Job, error) {
	allowed := s.store.AllowRequest(userID, s.maxRequests, s.window)
	if allowed {
		return true, nil, nil
	}

	if !s.queueEnabled || s.jobQueue == nil {
		return false, nil, errors.New("rate limit exceeded")
	}

	s.store.IncrementQueued(userID)

	job, err := s.jobQueue.Enqueue(userID, payload)
	if err != nil {
		return false, nil, err
	}

	return false, job, nil
}

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

func (s *RateLimiterService) GetJob(jobID string) (*model.Job, bool) {
	if s.jobQueue == nil {
		return nil, false
	}
	return s.jobQueue.GetJob(jobID)
}

func (s *RateLimiterService) GetQueueStats() model.QueueStatsResponse {
	if s.jobQueue == nil {
		return model.QueueStatsResponse{}
	}
	return s.jobQueue.GetStats()
}

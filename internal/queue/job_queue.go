package queue

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"rate-limited/internal/model"
)

type ProcessorFunc func(job *model.Job) error

type JobQueue struct {
	mu         sync.RWMutex
	jobs       map[string]*model.Job
	queue      chan *model.Job
	processing int
	succeeded  int
	failed     int
	maxRetries int
	processor  ProcessorFunc
}

func NewJobQueue(size int, maxRetries int, workerCount int, processor ProcessorFunc) *JobQueue {
	jq := &JobQueue{
		jobs:       make(map[string]*model.Job),
		queue:      make(chan *model.Job, size),
		maxRetries: maxRetries,
		processor:  processor,
	}

	for i := 0; i < workerCount; i++ {
		go jq.worker()
	}

	return jq
}

func (q *JobQueue) Enqueue(userID string, payload interface{}) (*model.Job, error) {
	job := &model.Job{
		ID:         uuid.NewString(),
		UserID:     userID,
		Payload:    payload,
		Status:     model.JobStatusQueued,
		Retries:    0,
		MaxRetries: q.maxRetries,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	q.mu.Lock()
	q.jobs[job.ID] = job
	q.mu.Unlock()

	select {
	case q.queue <- job:
		return job, nil
	default:
		return nil, fmt.Errorf("queue is full")
	}
}

func (q *JobQueue) worker() {
	for job := range q.queue {
		q.mu.Lock()
		q.processing++
		job.Status = model.JobStatusProcessing
		job.UpdatedAt = time.Now()
		q.mu.Unlock()

		err := q.processor(job)

		q.mu.Lock()
		q.processing--

		if err == nil {
			job.Status = model.JobStatusSucceeded
			job.Error = ""
			job.UpdatedAt = time.Now()
			q.succeeded++
			q.mu.Unlock()
			continue
		}

		job.Retries++
		job.Error = err.Error()
		job.UpdatedAt = time.Now()

		if job.Retries <= job.MaxRetries {
			job.Status = model.JobStatusQueued
			q.mu.Unlock()

			time.Sleep(time.Duration(job.Retries) * time.Second)

			select {
			case q.queue <- job:
			default:
				q.mu.Lock()
				job.Status = model.JobStatusFailed
				job.Error = "queue full while retrying"
				job.UpdatedAt = time.Now()
				q.failed++
				q.mu.Unlock()
			}
			continue
		}

		job.Status = model.JobStatusFailed
		q.failed++
		q.mu.Unlock()
	}
}

func (q *JobQueue) GetJob(jobID string) (*model.Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	job, exists := q.jobs[jobID]
	return job, exists
}

func (q *JobQueue) GetStats() model.QueueStatsResponse {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return model.QueueStatsResponse{
		Queued:     len(q.queue),
		Processing: q.processing,
		Succeeded:  q.succeeded,
		Failed:     q.failed,
	}
}

package model

import "time"

type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusSucceeded  JobStatus = "succeeded"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Payload    interface{} `json:"payload"`
	Status     JobStatus   `json:"status"`
	Retries    int         `json:"retries"`
	MaxRetries int         `json:"max_retries"`
	Error      string      `json:"error,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type QueueStatsResponse struct {
	Queued     int `json:"queued"`
	Processing int `json:"processing"`
	Succeeded  int `json:"succeeded"`
	Failed     int `json:"failed"`
}

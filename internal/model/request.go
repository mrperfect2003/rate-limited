package model

type RequestPayload struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

type RequestResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserStat struct {
	UserID         string `json:"user_id"`
	TotalAccepted  int    `json:"total_accepted"`
	TotalRejected  int    `json:"total_rejected"`
	QueuedRequests int    `json:"queued_requests"`
}

type StatsResponse struct {
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalUsers int        `json:"total_users"`
	Stats      []UserStat `json:"stats"`
}

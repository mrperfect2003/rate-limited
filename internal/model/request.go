package model

// RequestPayload represents input for POST /request.
type RequestPayload struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

// RequestResponse represents success response for POST /request.
type RequestResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// ErrorResponse represents generic API error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserStat represents request statistics for a single user.
type UserStat struct {
	UserID        string `json:"user_id"`
	TotalRequests int    `json:"total_requests"`
}

// StatsResponse represents paginated response for GET /stats.
type StatsResponse struct {
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalUsers int        `json:"total_users"`
	Stats      []UserStat `json:"stats"`
}

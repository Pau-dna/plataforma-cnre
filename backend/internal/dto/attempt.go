package dto

import "time"

// UpdateEvaluationAttemptRequest DTO for updating evaluation attempts (PATCH)
type UpdateEvaluationAttemptRequest struct {
	Score       *float64   `json:"score,omitempty"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

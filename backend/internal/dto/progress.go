package dto

import "time"

// UpdateUserProgressRequest DTO for updating user progress (PATCH)
type UpdateUserProgressRequest struct {
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Score       *int       `json:"score,omitempty"`
	Attempts    *int       `json:"attempts,omitempty"`
}

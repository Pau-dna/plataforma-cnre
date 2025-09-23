package dto

import "time"

// UpdateUserProgressRequest DTO for updating user progress (PATCH)
type UpdateUserProgressRequest struct {
	Progress    *float64   `json:"progress,omitempty"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	TimeSpent   *int       `json:"time_spent,omitempty"`
}

package dto

import "time"

// UpdateEvaluationAttemptRequest DTO for updating evaluation attempts (PATCH)
type UpdateEvaluationAttemptRequest struct {
	Score       *int       `json:"score,omitempty"`
	TotalPoints *int       `json:"total_points,omitempty"`
	Passed      *bool      `json:"passed,omitempty"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	TimeSpent   *int       `json:"time_spent,omitempty"`
}

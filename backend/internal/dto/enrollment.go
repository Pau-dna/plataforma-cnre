package dto

import "time"

// UpdateEnrollmentRequest DTO for updating enrollments (PATCH)
type UpdateEnrollmentRequest struct {
	Progress    *float64   `json:"progress,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

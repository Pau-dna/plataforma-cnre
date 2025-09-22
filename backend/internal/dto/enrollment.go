package dto

// UpdateEnrollmentRequest DTO for updating enrollments (PATCH)
type UpdateEnrollmentRequest struct {
	Progress   *float64 `json:"progress,omitempty"`
	IsComplete *bool    `json:"is_complete,omitempty"`
}
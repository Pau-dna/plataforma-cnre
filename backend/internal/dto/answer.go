package dto

// UpdateAnswerRequest DTO for updating answers (PATCH)
type UpdateAnswerRequest struct {
	Text      *string `json:"text,omitempty"`
	IsCorrect *bool   `json:"is_correct,omitempty"`
	Order     *int    `json:"order,omitempty"`
}
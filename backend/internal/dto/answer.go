package dto

// CreateAnswerRequest DTO for creating answers
type CreateAnswerRequest struct {
	Text       string `json:"text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
	Order      int    `json:"order" binding:"required"`
	QuestionID uint   `json:"question_id" binding:"required"`
}

// UpdateAnswerRequest DTO for updating answers (PUT)
type UpdateAnswerRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
	Order     int    `json:"order"`
}

// UpdateAnswerPatchRequest DTO for partially updating answers (PATCH)
type UpdateAnswerPatchRequest struct {
	Text      *string `json:"text,omitempty"`
	IsCorrect *bool   `json:"is_correct,omitempty"`
	Order     *int    `json:"order,omitempty"`
}

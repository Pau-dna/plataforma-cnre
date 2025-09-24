package dto

import "github.com/imlargo/go-api-template/internal/enums"

// UpdateEvaluationRequest DTO for updating evaluations (PATCH)
type UpdateEvaluationRequest struct {
	Order              *int               `json:"order,omitempty"`
	Title              *string            `json:"title,omitempty"`
	Description        *string            `json:"description,omitempty"`
	Type               *enums.ContentType `json:"type,omitempty"`
	QuestionCount      *int               `json:"question_count,omitempty"`
	AnswerOptionsCount *int               `json:"answer_options_count,omitempty"`
	PassingScore       *float64           `json:"passing_score,omitempty"`
	MaxAttempts        *int               `json:"max_attempts,omitempty"`
	TimeLimit          *int               `json:"time_limit,omitempty"`
}

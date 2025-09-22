package dto

import "github.com/imlargo/go-api-template/internal/enums"

// UpdateQuestionRequest DTO for updating questions (PATCH)
type UpdateQuestionRequest struct {
	Text        *string              `json:"text,omitempty"`
	Type        *enums.QuestionType  `json:"type,omitempty"`
	Explanation *string              `json:"explanation,omitempty"`
	Points      *float64             `json:"points,omitempty"`
}
package dto

import "github.com/imlargo/go-api-template/internal/enums"

// CreateQuestionRequest DTO for creating questions
type CreateQuestionRequest struct {
	Text         string             `json:"text" binding:"required"`
	Type         enums.QuestionType `json:"type" binding:"required"`
	Explanation  string             `json:"explanation"`
	Points       int                `json:"points" binding:"required"`
	EvaluationID uint               `json:"evaluation_id" binding:"required"`
}

// UpdateQuestionRequest DTO for updating questions (PUT)
type UpdateQuestionRequest struct {
	Text        string             `json:"text"`
	Type        enums.QuestionType `json:"type"`
	Explanation string             `json:"explanation"`
	Points      int                `json:"points"`
}

// UpdateQuestionPatchRequest DTO for partially updating questions (PATCH)
type UpdateQuestionPatchRequest struct {
	Text        *string             `json:"text,omitempty"`
	Type        *enums.QuestionType `json:"type,omitempty"`
	Explanation *string             `json:"explanation,omitempty"`
	Points      *int                `json:"points,omitempty"`
}

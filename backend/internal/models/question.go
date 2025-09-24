package models

import "github.com/imlargo/go-api-template/internal/enums"

// Question - modelo de pregunta
type Question struct {
	BaseModel
	Text         string             `json:"text" gorm:"type:text;not null"`
	Type         enums.QuestionType `json:"type" gorm:"not null"`
	Explanation  string             `json:"explanation" gorm:"type:text"`
	Points       int                `json:"points" gorm:"not null;default:1"`
	EvaluationID uint               `json:"evaluation_id" gorm:"not null;index"`

	// Relaciones
	Evaluation *Evaluation `json:"evaluation" gorm:"foreignKey:EvaluationID;constraint:OnDelete:CASCADE"`
	Answers    []*Answer   `json:"answers" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

func (Question) TableName() string {
	return "questions"
}

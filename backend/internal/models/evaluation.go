package models

import "github.com/imlargo/go-api-template/internal/enums"

// Evaluation - modelo de evaluación (quizzes, exámenes)
type Evaluation struct {
	BaseModel
	Order         int               `json:"order" gorm:"not null"`
	Title         string            `json:"title" gorm:"not null"`
	Description   *string           `json:"description" gorm:"type:text"`
	Type          enums.ContentType `json:"type" gorm:"not null;default:'evaluation'"`
	QuestionCount int               `json:"question_count" gorm:"not null"`
	PassingScore  int               `json:"passing_score" gorm:"not null"`
	MaxAttempts   *int              `json:"max_attempts"`
	TimeLimit     *int              `json:"time_limit"` // en minutos
	ModuleID      uint              `json:"module_id" gorm:"not null"`

	// Relaciones
	Module             Module              `json:"module" gorm:"foreignKey:ModuleID"`
	Questions          []Question          `json:"questions" gorm:"foreignKey:EvaluationID"`
	EvaluationAttempts []EvaluationAttempt `json:"evaluation_attempts" gorm:"foreignKey:EvaluationID"`
	UserProgress       []UserProgress      `json:"user_progress" gorm:"foreignKey:ContentID"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// AttemptAnswer - estructura para las respuestas de un intento
type AttemptAnswer struct {
	QuestionID        uint   `json:"question_id"`
	SelectedAnswerIDs []uint `json:"selected_answer_ids"`
	IsCorrect         bool   `json:"is_correct"`
	Points            int    `json:"points"`
}

// AttemptAnswers - slice personalizado para manejar JSON
type AttemptAnswers []AttemptAnswer

// Implementar driver.Valuer para poder guardar en la base de datos
func (a AttemptAnswers) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Implementar sql.Scanner para poder leer desde la base de datos
func (a *AttemptAnswers) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte into AttemptAnswers")
	}

	return json.Unmarshal(bytes, a)
}

// EvaluationAttempt - modelo de intento de evaluación
type EvaluationAttempt struct {
	BaseModel
	UserID       uint           `json:"user_id" gorm:"not null"`
	EvaluationID uint           `json:"evaluation_id" gorm:"not null"`
	Answers      AttemptAnswers `json:"answers" gorm:"type:json"`
	Score        int            `json:"score" gorm:"not null;default:0"`
	TotalPoints  int            `json:"total_points" gorm:"not null"`
	Passed       bool           `json:"passed" gorm:"not null;default:false"`
	StartedAt    time.Time      `json:"started_at" gorm:"not null"`
	SubmittedAt  *time.Time     `json:"submitted_at"`
	TimeSpent    *int           `json:"time_spent"` // en minutos

	// Relaciones
	User       User       `json:"user" gorm:"foreignKey:UserID"`
	Evaluation Evaluation `json:"evaluation" gorm:"foreignKey:EvaluationID"`
}

func (EvaluationAttempt) TableName() string {
	return "evaluation_attempts"
}

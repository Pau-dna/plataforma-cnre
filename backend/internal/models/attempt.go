package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/imlargo/go-api-template/internal/enums"
)

// AttemptAnswerOption - opción de respuesta generada para una pregunta del intento
type AttemptAnswerOption struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// AttemptQuestion - pregunta generada para un intento específico
type AttemptQuestion struct {
	ID            uint                  `json:"id"`
	Text          string                `json:"text"`
	Type          enums.QuestionType    `json:"type"`
	Explanation   string                `json:"explanation"`
	Points        int                   `json:"points"`
	OriginalID    uint                  `json:"original_id"`    // ID de la pregunta original
	AnswerOptions []AttemptAnswerOption `json:"answer_options"` // Opciones de respuesta generadas
}

// AttemptQuestions - slice personalizado para manejar JSON
type AttemptQuestions []AttemptQuestion

// Implementar driver.Valuer para poder guardar en la base de datos
func (aq AttemptQuestions) Value() (driver.Value, error) {
	return json.Marshal(aq)
}

// Implementar sql.Scanner para poder leer desde la base de datos
func (aq *AttemptQuestions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte into AttemptQuestions")
	}

	return json.Unmarshal(bytes, aq)
}

// AttemptAnswer - estructura para las respuestas de un intento
type AttemptAnswer struct {
	AttemptQuestionID uint   `json:"attempt_question_id"` // ID de la pregunta generada del intento
	SelectedOptionIDs []uint `json:"selected_option_ids"` // IDs de las opciones seleccionadas
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
	UserID       uint             `json:"user_id" gorm:"not null;index;index:idx_eval_attempts_user_eval,priority:1"`
	EvaluationID uint             `json:"evaluation_id" gorm:"not null;index;index:idx_eval_attempts_user_eval,priority:2"`
	Questions    AttemptQuestions `json:"questions" gorm:"type:json"` // Preguntas generadas para este intento
	Answers      AttemptAnswers   `json:"answers" gorm:"type:json"`   // Respuestas del usuario
	Score        int              `json:"score" gorm:"not null;default:0"`
	TotalPoints  int              `json:"total_points" gorm:"not null"`
	Passed       bool             `json:"passed" gorm:"not null;default:false"`
	StartedAt    time.Time        `json:"started_at" gorm:"not null"`
	SubmittedAt  time.Time        `json:"submitted_at"`
	TimeSpent    int              `json:"time_spent"` // en minutos

	// Relaciones
	User       *User       `json:"user" gorm:"foreignKey:UserID"`
	Evaluation *Evaluation `json:"evaluation" gorm:"foreignKey:EvaluationID"`
}

func (EvaluationAttempt) TableName() string {
	return "evaluation_attempts"
}

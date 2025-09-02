package models

import (
	"time"
)

type Course struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string
}

type Module struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name  string
	Order int
}

type Content struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string
	Url  string
}

type Evaluation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	QuestionCount int
}

type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      QuestionType
}

type QuestionType string

const (
	QuestionSingle   QuestionType = "single"
	QuestionMultiple QuestionType = "multiple"
)

type Answer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Correct   bool
}

type EnrollMent struct {
	UserID   int
	CourseID uint
}

type Grade struct {
	EvaluationID uint
	Grade        float64
}

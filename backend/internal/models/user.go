package models

import "github.com/imlargo/go-api-template/internal/enums"

type User struct {
	BaseModel

	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	FirstName string         `json:"firstName" gorm:"not null"`
	LastName  string         `json:"lastName" gorm:"not null"`
	AvatarURL *string        `json:"avatarUrl" gorm:"column:avatar_url"`
	Role      enums.UserRole `json:"role" gorm:"not null;default:'student'"`

	// Relaciones
	Enrollments        []Enrollment        `json:"enrollments" gorm:"foreignKey:UserID"`
	EvaluationAttempts []EvaluationAttempt `json:"evaluation_attempts" gorm:"foreignKey:UserID"`
	UserProgress       []UserProgress      `json:"user_progress" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

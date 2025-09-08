package models

import "time"

// Enrollment - modelo de inscripción a curso
type Enrollment struct {
	BaseModel
	UserID      uint       `json:"user_id" gorm:"not null"`
	CourseID    uint       `json:"course_id" gorm:"not null"`
	EnrolledAt  time.Time  `json:"enrolled_at" gorm:"not null"`
	CompletedAt *time.Time `json:"completed_at"`
	Progress    float64    `json:"progress" gorm:"not null;default:0.0"` // porcentaje 0-100

	// Relaciones
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Course Course `json:"course" gorm:"foreignKey:CourseID"`

	// Índice compuesto para evitar inscripciones duplicadas
	// Se define en el método de migración: gorm:"uniqueIndex:idx_user_course"
}

func (Enrollment) TableName() string {
	return "enrollments"
}

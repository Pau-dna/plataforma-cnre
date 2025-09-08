package models

import "time"

// UserProgress - modelo de progreso del usuario
type UserProgress struct {
	BaseModel
	UserID      uint      `json:"user_id" gorm:"not null"`
	CourseID    uint      `json:"course_id" gorm:"not null"`
	ModuleID    uint      `json:"module_id" gorm:"not null"`
	ContentID   uint      `json:"content_id" gorm:"not null"` // puede ser Content o Evaluation
	CompletedAt time.Time `json:"completed_at"`
	Score       int       `json:"score"`
	Attempts    int       `json:"attempts" gorm:"not null;default:0"`

	// Relaciones
	User    *User    `json:"user" gorm:"foreignKey:UserID"`
	Course  *Course  `json:"course" gorm:"foreignKey:CourseID"`
	Module  *Module  `json:"module" gorm:"foreignKey:ModuleID"`
	Content *Content `json:"content" gorm:"foreignKey:ContentID"`
}

func (UserProgress) TableName() string {
	return "user_progress"
}

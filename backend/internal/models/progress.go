package models

import "time"

// UserProgress - modelo de progreso del usuario
type UserProgress struct {
	BaseModel
	UserID      uint      `json:"user_id" gorm:"not null;index;index:idx_user_progress_user_course,priority:1;index:idx_user_progress_user_module,priority:1;index:idx_user_progress_user_content,priority:1"`
	CourseID    uint      `json:"course_id" gorm:"not null;index;index:idx_user_progress_user_course,priority:2"`
	ModuleID    uint      `json:"module_id" gorm:"not null;index;index:idx_user_progress_user_module,priority:2"`
	ContentID   uint      `json:"content_id" gorm:"not null;index;index:idx_user_progress_user_content,priority:2"` // puede ser Content o Evaluation
	CompletedAt time.Time `json:"completed_at"`
	Score       int       `json:"score"`
	Attempts    int       `json:"attempts" gorm:"not null;default:0"`

	// Relaciones
	User    *User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Course  *Course  `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Module  *Module  `json:"module" gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE"`
	Content *Content `json:"content" gorm:"foreignKey:ContentID;constraint:OnDelete:CASCADE"`
}

func (UserProgress) TableName() string {
	return "user_progress"
}

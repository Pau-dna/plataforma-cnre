package models

// Module - modelo de módulo
type Module struct {
	BaseModel
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Order       int    `json:"order" gorm:"not null;index:idx_modules_course_order,priority:2"`
	CourseID    uint   `json:"course_id" gorm:"not null;index;index:idx_modules_course_order,priority:1"`

	// Relaciones
	Course      *Course       `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Contents    []*Content    `json:"contents" gorm:"foreignKey:ModuleID"`
	Evaluations []*Evaluation `json:"evaluations" gorm:"foreignKey:ModuleID"`
}

func (Module) TableName() string {
	return "modules"
}

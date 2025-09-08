package models

// Module - modelo de módulo
type Module struct {
	BaseModel
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Order       int    `json:"order" gorm:"not null"`
	CourseID    uint   `json:"course_id" gorm:"not null"`

	// Relaciones
	Course      *Course       `json:"course" gorm:"foreignKey:CourseID"`
	Contents    []*Content    `json:"contents" gorm:"foreignKey:ModuleID"`
	Evaluations []*Evaluation `json:"evaluations" gorm:"foreignKey:ModuleID"`
}

func (Module) TableName() string {
	return "modules"
}

package models

// Course - modelo de curso
type Course struct {
	BaseModel
	Title            string  `json:"title" gorm:"not null"`
	Description      string  `json:"description" gorm:"type:text"`
	ShortDescription *string `json:"shortDescription" gorm:"column:short_description;type:text"`
	ImageURL         *string `json:"imageUrl" gorm:"column:image_url"`
	StudentCount     *int    `json:"student_count" gorm:"column:student_count"`
	ModuleCount      *int    `json:"module_count" gorm:"column:module_count"`

	// Relaciones
	Modules      []Module       `json:"modules" gorm:"foreignKey:CourseID"`
	Enrollments  []Enrollment   `json:"enrollments" gorm:"foreignKey:CourseID"`
	UserProgress []UserProgress `json:"user_progress" gorm:"foreignKey:CourseID"`
}

func (Course) TableName() string {
	return "courses"
}

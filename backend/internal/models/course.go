package models

// Course - modelo de curso
type Course struct {
	BaseModel
	Title            string `json:"title" gorm:"not null"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	StudentCount     int    `json:"student_count"`
	ModuleCount      int    `json:"module_count"`

	// Relaciones
	Modules     []*Module     `json:"modules" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Enrollments []*Enrollment `json:"enrollments" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}

func (Course) TableName() string {
	return "courses"
}

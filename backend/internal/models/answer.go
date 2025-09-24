package models

// Answer - modelo de respuesta para las opciones de pregunta
type Answer struct {
	BaseModel
	Text       string `json:"text" gorm:"type:text;not null"`
	IsCorrect  bool   `json:"is_correct" gorm:"column:is_correct;not null;default:false"`
	Order      int    `json:"order" gorm:"not null;index:idx_answers_question_order,priority:2"`
	QuestionID uint   `json:"question_id" gorm:"not null;index;index:idx_answers_question_order,priority:1"`

	// Relaciones
	Question *Question `json:"question" gorm:"foreignKey:QuestionID"`
}

func (Answer) TableName() string {
	return "answers"
}

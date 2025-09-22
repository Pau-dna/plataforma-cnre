package models

import "github.com/imlargo/go-api-template/internal/enums"

// Content - modelo de contenido (lecciones, videos, lecturas)
type Content struct {
	BaseModel
	Order       int               `json:"order" gorm:"not null"`
	Title       string            `json:"title" gorm:"not null"`
	Description string            `json:"description" gorm:"type:text"`
	Type        enums.ContentType `json:"type" gorm:"not null;default:'content'"`
	Body        string            `json:"body" gorm:"type:text;not null"`
	MediaURL    string            `json:"media_url" gorm:"column:media_url"`
	ModuleID    uint              `json:"module_id" gorm:"not null"`

	// Relaciones
	Module       *Module         `json:"module" gorm:"foreignKey:ModuleID"`
	UserProgress []*UserProgress `json:"user_progress" gorm:"foreignKey:ContentID"`
}

func (Content) TableName() string {
	return "contents"
}

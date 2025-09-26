package models

import "github.com/imlargo/go-api-template/internal/enums"

// Content - modelo de contenido (lecciones, videos, lecturas)
type Content struct {
	BaseModel
	Order       int               `json:"order" gorm:"not null;index:idx_contents_module_order,priority:2"`
	Title       string            `json:"title" gorm:"not null"`
	Description string            `json:"description" gorm:"type:text"`
	Type        enums.ContentType `json:"type" gorm:"not null;default:'content'"`
	Body        string            `json:"body" gorm:"type:text;not null"`
	MediaURL    string            `json:"media_url" gorm:"column:media_url"`
	ModuleID    uint              `json:"module_id" gorm:"not null;index;index:idx_contents_module_order,priority:1"`

	// Relaciones
	Module       *Module         `json:"module" gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE"`
	UserProgress []*UserProgress `json:"user_progress" gorm:"foreignKey:ContentID;constraint:OnDelete:CASCADE"`
}

func (Content) TableName() string {
	return "contents"
}

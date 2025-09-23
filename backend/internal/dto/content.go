package dto

import "github.com/imlargo/go-api-template/internal/enums"

// UpdateContentRequest DTO for updating content (PATCH)
type UpdateContentRequest struct {
	Order       *int               `json:"order,omitempty"`
	Title       *string            `json:"title,omitempty"`
	Description *string            `json:"description,omitempty"`
	Type        *enums.ContentType `json:"type,omitempty"`
	Body        *string            `json:"body,omitempty"`
	MediaURL    *string            `json:"media_url,omitempty"`
}

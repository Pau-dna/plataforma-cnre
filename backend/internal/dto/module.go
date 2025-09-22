package dto

// UpdateModuleRequest DTO for updating modules (PATCH)
type UpdateModuleRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Order       *int    `json:"order,omitempty"`
}
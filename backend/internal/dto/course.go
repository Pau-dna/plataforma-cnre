package dto

// UpdateCourseRequest DTO for updating courses (PATCH)
type UpdateCourseRequest struct {
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	ShortDescription *string `json:"short_description,omitempty"`
	ImageURL         *string `json:"image_url,omitempty"`
	StudentCount     *int    `json:"student_count,omitempty"`
	ModuleCount      *int    `json:"module_count,omitempty"`
}

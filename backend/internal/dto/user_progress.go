package dto

// ModuleProgressDetail represents the progress details for a specific module
type ModuleProgressDetail struct {
	ModuleID     uint    `json:"module_id"`
	ModuleTitle  string  `json:"module_title"`
	Percentage   float64 `json:"percentage"`
	IsCompleted  bool    `json:"is_completed"`
}

// CourseProgressSummary represents a comprehensive course progress summary
type CourseProgressSummary struct {
	CourseID         uint                   `json:"course_id"`
	CourseTitle      string                 `json:"course_title"`
	TotalPercentage  float64                `json:"total_percentage"`
	IsCompleted      bool                   `json:"is_completed"`
	ModulesProgress  []ModuleProgressDetail `json:"modules_progress"`
}

// ContentItemDetail represents the progress details for a specific content item
type ContentItemDetail struct {
	ItemID      uint    `json:"item_id"`
	ItemTitle   string  `json:"item_title"`
	ItemType    string  `json:"item_type"` // "content" or "evaluation"
	IsCompleted bool    `json:"is_completed"`
	CompletedAt *string `json:"completed_at,omitempty"`
	Score       *int    `json:"score,omitempty"`
	Order       int     `json:"order"`
}

// ModuleProgressSummary represents a comprehensive module progress summary
type ModuleProgressSummary struct {
	ModuleID        uint                `json:"module_id"`
	ModuleTitle     string              `json:"module_title"`
	CourseID        uint                `json:"course_id"`
	CourseTitle     string              `json:"course_title"`
	TotalPercentage float64             `json:"total_percentage"`
	IsCompleted     bool                `json:"is_completed"`
	ContentItems    []ContentItemDetail `json:"content_items"`
}
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
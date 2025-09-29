package dto

// CourseKPIResponse DTO for course KPI dashboard metrics
type CourseKPIResponse struct {
	CourseID         uint    `json:"course_id"`
	CourseTitle      string  `json:"course_title"`
	StudentCount     int     `json:"student_count"`
	CompletionRate   float64 `json:"completion_rate"`   // percentage 0-100
	AverageProgress  float64 `json:"average_progress"`  // percentage 0-100
}
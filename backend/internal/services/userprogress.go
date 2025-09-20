package services

import (
	"fmt"
	"time"

	"github.com/imlargo/go-api-template/internal/models"
)

type UserProgressService interface {
	MarkContentComplete(userID, courseID, moduleID, contentID uint) (*models.UserProgress, error)
	MarkContentIncomplete(userID, courseID, moduleID, contentID uint) error
	GetUserProgress(userID, courseID uint) ([]*models.UserProgress, error)
	GetUserModuleProgress(userID, moduleID uint) ([]*models.UserProgress, error)
	CalculateCourseProgress(userID, courseID uint) (float64, error)
	CalculateModuleProgress(userID, moduleID uint) (float64, error)
	GetUserProgressForContent(userID, contentID uint) (*models.UserProgress, error)
}

type userProgressService struct {
	*Service
	enrollmentService EnrollmentService
}

func NewUserProgressService(service *Service, enrollmentService EnrollmentService) UserProgressService {
	return &userProgressService{
		Service:           service,
		enrollmentService: enrollmentService,
	}
}

func (s *userProgressService) MarkContentComplete(userID, courseID, moduleID, contentID uint) (*models.UserProgress, error) {
	// Verify user is enrolled in the course
	_, err := s.enrollmentService.GetUserCourseEnrollment(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("user not enrolled in course: %w", err)
	}

	// Check if progress already exists
	existing, _ := s.GetUserProgressForContent(userID, contentID)
	if existing != nil {
		// Update existing progress
		existing.CompletedAt = time.Now()
		existing.Attempts++
		
		if err := s.store.UserProgresss.Update(existing); err != nil {
			return nil, fmt.Errorf("failed to update progress: %w", err)
		}
		
		// Update course progress
		s.updateCourseProgress(userID, courseID)
		
		return existing, nil
	}

	// Create new progress record
	progress := &models.UserProgress{
		UserID:      userID,
		CourseID:    courseID,
		ModuleID:    moduleID,
		ContentID:   contentID,
		CompletedAt: time.Now(),
		Attempts:    1,
	}

	if err := s.store.UserProgresss.Create(progress); err != nil {
		return nil, fmt.Errorf("failed to create progress: %w", err)
	}

	// Update course progress
	s.updateCourseProgress(userID, courseID)

	return progress, nil
}

func (s *userProgressService) MarkContentIncomplete(userID, courseID, moduleID, contentID uint) error {
	// Find existing progress
	existing, err := s.GetUserProgressForContent(userID, contentID)
	if err != nil {
		return fmt.Errorf("progress not found: %w", err)
	}

	// Reset completion
	existing.CompletedAt = time.Time{}
	
	if err := s.store.UserProgresss.Update(existing); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	// Update course progress
	s.updateCourseProgress(userID, courseID)

	return nil
}

func (s *userProgressService) GetUserProgress(userID, courseID uint) ([]*models.UserProgress, error) {
	// This would require a repository method to filter by user ID and course ID
	// For now, we'll implement a basic version
	allProgress, err := s.store.UserProgresss.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Filter by user ID and course ID
	var userProgress []*models.UserProgress
	for _, progress := range allProgress {
		if progress.UserID == userID && progress.CourseID == courseID {
			userProgress = append(userProgress, progress)
		}
	}

	return userProgress, nil
}

func (s *userProgressService) GetUserModuleProgress(userID, moduleID uint) ([]*models.UserProgress, error) {
	// This would require a repository method to filter by user ID and module ID
	// For now, we'll implement a basic version
	allProgress, err := s.store.UserProgresss.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Filter by user ID and module ID
	var moduleProgress []*models.UserProgress
	for _, progress := range allProgress {
		if progress.UserID == userID && progress.ModuleID == moduleID {
			moduleProgress = append(moduleProgress, progress)
		}
	}

	return moduleProgress, nil
}

func (s *userProgressService) CalculateCourseProgress(userID, courseID uint) (float64, error) {
	// Get all modules for the course
	modules, err := s.store.Modules.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to get modules: %w", err)
	}

	var courseModules []*models.Module
	for _, module := range modules {
		if module.CourseID == courseID {
			courseModules = append(courseModules, module)
		}
	}

	if len(courseModules) == 0 {
		return 0, nil
	}

	totalModules := len(courseModules)
	completedModules := 0

	// Calculate progress for each module
	for _, module := range courseModules {
		moduleProgress, err := s.CalculateModuleProgress(userID, module.ID)
		if err != nil {
			continue
		}
		
		// Consider module complete if 100% progress
		if moduleProgress >= 100.0 {
			completedModules++
		}
	}

	return float64(completedModules) / float64(totalModules) * 100.0, nil
}

func (s *userProgressService) CalculateModuleProgress(userID, moduleID uint) (float64, error) {
	// Get all content for the module
	contents, err := s.store.Contents.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to get contents: %w", err)
	}

	// Get all evaluations for the module
	evaluations, err := s.store.Evaluations.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to get evaluations: %w", err)
	}

	// Count total items in module
	totalItems := 0
	for _, content := range contents {
		if content.ModuleID == moduleID {
			totalItems++
		}
	}
	for _, evaluation := range evaluations {
		if evaluation.ModuleID == moduleID {
			totalItems++
		}
	}

	if totalItems == 0 {
		return 100.0, nil // No content means 100% complete
	}

	// Get user progress for this module
	userProgress, err := s.GetUserModuleProgress(userID, moduleID)
	if err != nil {
		return 0, err
	}

	// Count completed items
	completedItems := 0
	for _, progress := range userProgress {
		if !progress.CompletedAt.IsZero() {
			completedItems++
		}
	}

	return float64(completedItems) / float64(totalItems) * 100.0, nil
}

func (s *userProgressService) GetUserProgressForContent(userID, contentID uint) (*models.UserProgress, error) {
	// This would require a repository method to filter by user ID and content ID
	allProgress, err := s.store.UserProgresss.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Find progress for specific content
	for _, progress := range allProgress {
		if progress.UserID == userID && progress.ContentID == contentID {
			return progress, nil
		}
	}

	return nil, fmt.Errorf("progress not found")
}

func (s *userProgressService) updateCourseProgress(userID, courseID uint) error {
	// Calculate new course progress
	progress, err := s.CalculateCourseProgress(userID, courseID)
	if err != nil {
		s.logger.Warnf("Failed to calculate course progress for user %d, course %d: %v", userID, courseID, err)
		return err
	}

	// Update enrollment progress
	if s.enrollmentService != nil {
		err = s.enrollmentService.UpdateProgress(userID, courseID, progress)
		if err != nil {
			s.logger.Warnf("Failed to update enrollment progress: %v", err)
		}
	}

	return nil
}
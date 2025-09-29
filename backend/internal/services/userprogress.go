package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type UserProgressService interface {
	MarkContentComplete(userID, courseID, moduleID, contentID uint) (*models.UserProgress, error)
	MarkContentIncomplete(userID, courseID, moduleID, contentID uint) error
	UpdateUserProgressPatch(id uint, data map[string]interface{}) (*models.UserProgress, error)
	GetUserProgress(userID, courseID uint) ([]*models.UserProgress, error)
	GetUserModuleProgress(userID, moduleID uint) ([]*models.UserProgress, error)
	CalculateCourseProgress(userID, courseID uint) (float64, error)
	CalculateModuleProgress(userID, moduleID uint) (float64, error)
	GetUserProgressForContent(userID, contentID uint) (*models.UserProgress, error)
	HasUserPassedEvaluation(userID, evaluationID uint) (bool, error)
	UpdateCourseProgress(userID, courseID uint) error
	GetComprehensiveCourseProgress(userID, courseID uint) (*dto.CourseProgressSummary, error)
	GetModuleContentProgress(userID, moduleID uint) ([]*dto.ContentProgressResponse, error)
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
		return nil, fmt.Errorf("el usuario no está inscrito en el curso: %w", err)
	}

	// Check if progress already exists
	existing, _ := s.GetUserProgressForContent(userID, contentID)
	if existing != nil {
		// Update existing progress
		existing.CompletedAt = time.Now()
		existing.Attempts++

		if err := s.store.UserProgresss.Update(existing); err != nil {
			return nil, fmt.Errorf("error al actualizar el progreso: %w", err)
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
		return nil, fmt.Errorf("error al crear el progreso: %w", err)
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
		return fmt.Errorf("error al actualizar el progreso: %w", err)
	}

	// Update course progress
	s.updateCourseProgress(userID, courseID)

	return nil
}

func (s *userProgressService) UpdateUserProgressPatch(progressID uint, data map[string]interface{}) (*models.UserProgress, error) {
	if progressID == 0 {
		return nil, errors.New("progress ID cannot be zero")
	}

	var progress dto.UpdateUserProgressRequest
	if err := utils.MapToStructStrict(data, &progress); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.UserProgresss.Patch(progressID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.UserProgresss.Get(progressID)
	if err != nil {
		return nil, errors.New("progreso no encontrado")
	}

	return updated, nil
}

func (s *userProgressService) GetUserProgress(userID, courseID uint) ([]*models.UserProgress, error) {
	// Use the new repository method to filter by user ID and course ID at database level
	userProgress, err := s.store.UserProgresss.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el progreso: %w", err)
	}

	return userProgress, nil
}

func (s *userProgressService) GetUserModuleProgress(userID, moduleID uint) ([]*models.UserProgress, error) {
	// Use the new repository method to filter by user ID and module ID at database level
	moduleProgress, err := s.store.UserProgresss.GetByUserAndModule(userID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el progreso: %w", err)
	}

	return moduleProgress, nil
}

func (s *userProgressService) CalculateCourseProgress(userID, courseID uint) (float64, error) {
	// Use optimized query to get modules for the specific course
	courseModules, err := s.store.Modules.GetByCourseID(courseID)
	if err != nil {
		return 0, fmt.Errorf("error al obtener los módulos: %w", err)
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
	// Use optimized queries to get content and evaluations for the specific module
	contents, err := s.store.Contents.GetByModuleID(moduleID)
	if err != nil {
		return 0, fmt.Errorf("error al obtener los contenidos: %w", err)
	}

	evaluations, err := s.store.Evaluations.GetByModuleID(moduleID)
	if err != nil {
		return 0, fmt.Errorf("error al obtener las evaluaciones: %w", err)
	}

	// Count total items in module
	moduleContents := []*models.Content{}
	moduleEvaluations := []*models.Evaluation{}
	
	for _, content := range contents {
		if content.ModuleID == moduleID {
			moduleContents = append(moduleContents, content)
		}
	}
	for _, evaluation := range evaluations {
		if evaluation.ModuleID == moduleID {
			moduleEvaluations = append(moduleEvaluations, evaluation)
		}
	}

	totalItems := len(moduleContents) + len(moduleEvaluations)
	if totalItems == 0 {
		return 100.0, nil // No content means 100% complete
	}

	// Count completed items
	completedItems := 0

	// Check completed content (user progress records)
	userProgress, err := s.GetUserModuleProgress(userID, moduleID)
	if err != nil {
		return 0, err
	}

	// Create a map to quickly lookup completed content
	completedContentMap := make(map[uint]bool)
	for _, progress := range userProgress {
		if !progress.CompletedAt.IsZero() {
			completedContentMap[progress.ContentID] = true
		}
	}

	// Count completed content items
	for _, content := range moduleContents {
		if completedContentMap[content.ID] {
			completedItems++
		}
	}

	// Count passed evaluations (not just completed ones)
	for _, evaluation := range moduleEvaluations {
		hasPassed, err := s.HasUserPassedEvaluation(userID, evaluation.ID)
		if err != nil {
			s.logger.Warnf("Failed to check if user %d passed evaluation %d: %v", userID, evaluation.ID, err)
			continue
		}
		if hasPassed {
			completedItems++
		}
	}

	return float64(completedItems) / float64(totalItems) * 100.0, nil
}

func (s *userProgressService) GetUserProgressForContent(userID, contentID uint) (*models.UserProgress, error) {
	// Use the new repository method to filter by user ID and content ID at database level
	progress, err := s.store.UserProgresss.GetByUserAndContent(userID, contentID)
	if err != nil {
		return nil, fmt.Errorf("progress not found: %w", err)
	}

	return progress, nil
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

// UpdateCourseProgress - método público para actualizar el progreso del curso
func (s *userProgressService) UpdateCourseProgress(userID, courseID uint) error {
	return s.updateCourseProgress(userID, courseID)
}

func (s *userProgressService) HasUserPassedEvaluation(userID, evaluationID uint) (bool, error) {
	// Get all attempts for this user and evaluation
	attempts, err := s.store.EvaluationAttempts.GetAll()
	if err != nil {
		return false, fmt.Errorf("failed to get evaluation attempts: %w", err)
	}

	// Filter attempts for the specific user and evaluation
	var userAttempts []*models.EvaluationAttempt
	for _, attempt := range attempts {
		if attempt.UserID == userID && attempt.EvaluationID == evaluationID {
			userAttempts = append(userAttempts, attempt)
		}
	}

	// Check if any attempt was passed
	for _, attempt := range userAttempts {
		if attempt.Passed && attempt.SubmittedAt != nil && !attempt.SubmittedAt.IsZero() {
			return true, nil
		}
	}

	return false, nil
}

func (s *userProgressService) GetComprehensiveCourseProgress(userID, courseID uint) (*dto.CourseProgressSummary, error) {
	// Delegate to repository layer which handles the complex SQL query and data processing
	return s.store.UserProgresss.GetCourseProgressSummary(userID, courseID)
}

func (s *userProgressService) GetModuleContentProgress(userID, moduleID uint) ([]*dto.ContentProgressResponse, error) {
	// Use the optimized repository method to get content progress in a single query
	results, err := s.store.Contents.GetContentProgressByModule(userID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el progreso de contenidos del módulo: %w", err)
	}

	// Convert repository results to DTO
	contentProgress := make([]*dto.ContentProgressResponse, len(results))
	for i, result := range results {
		contentProgress[i] = &dto.ContentProgressResponse{
			ID:        result.ID,
			Title:     result.Title,
			Completed: result.Completed,
		}
	}

	return contentProgress, nil
}

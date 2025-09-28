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
	GetComprehensiveCourseProgress(userID, courseID uint) (*dto.CourseProgressSummary, error)
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
	// Get course information
	course, err := s.store.Courses.Get(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el curso: %w", err)
	}

	// Get all modules for the course in a single query
	modules, err := s.store.Modules.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los módulos: %w", err)
	}

	if len(modules) == 0 {
		// No modules means course is 100% complete
		return &dto.CourseProgressSummary{
			CourseID:        course.ID,
			CourseTitle:     course.Title,
			TotalPercentage: 100.0,
			IsCompleted:     true,
			ModulesProgress: []dto.ModuleProgressDetail{},
		}, nil
	}

	// Extract module IDs for batch queries
	moduleIDs := make([]uint, len(modules))
	moduleMap := make(map[uint]*models.Module)
	for i, module := range modules {
		moduleIDs[i] = module.ID
		moduleMap[module.ID] = module
	}

	// Batch fetch all contents and evaluations for all modules
	allContents, err := s.batchGetContentsByModuleIDs(moduleIDs)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los contenidos: %w", err)
	}

	allEvaluations, err := s.batchGetEvaluationsByModuleIDs(moduleIDs)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las evaluaciones: %w", err)
	}

	// Batch fetch all user progress for the course
	allUserProgress, err := s.GetUserProgress(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener progreso del usuario: %w", err)
	}

	// Batch fetch all evaluation attempts for the user
	allAttempts, err := s.batchGetUserEvaluationAttempts(userID, moduleIDs)
	if err != nil {
		return nil, fmt.Errorf("error al obtener intentos de evaluación: %w", err)
	}

	// Create lookup maps for efficient processing
	contentsByModule := s.groupContentsByModule(allContents)
	evaluationsByModule := s.groupEvaluationsByModule(allEvaluations)
	progressByContent := s.groupProgressByContent(allUserProgress)
	attemptsByEvaluation := s.groupAttemptsByEvaluation(allAttempts)

	// Calculate progress for all modules efficiently
	var modulesProgress []dto.ModuleProgressDetail
	completedModules := 0

	for _, module := range modules {
		moduleContents := contentsByModule[module.ID]
		moduleEvaluations := evaluationsByModule[module.ID]

		totalItems := len(moduleContents) + len(moduleEvaluations)
		if totalItems == 0 {
			// No content means 100% complete
			modulesProgress = append(modulesProgress, dto.ModuleProgressDetail{
				ModuleID:    module.ID,
				ModuleTitle: module.Title,
				Percentage:  100.0,
				IsCompleted: true,
			})
			completedModules++
			continue
		}

		completedItems := 0

		// Count completed content items
		for _, content := range moduleContents {
			if progress, exists := progressByContent[content.ID]; exists && !progress.CompletedAt.IsZero() {
				completedItems++
			}
		}

		// Count passed evaluations
		for _, evaluation := range moduleEvaluations {
			if attempts, exists := attemptsByEvaluation[evaluation.ID]; exists {
				for _, attempt := range attempts {
					if attempt.Passed && attempt.SubmittedAt != nil && !attempt.SubmittedAt.IsZero() {
						completedItems++
						break // Only count once per evaluation
					}
				}
			}
		}

		modulePercentage := float64(completedItems) / float64(totalItems) * 100.0
		isCompleted := modulePercentage >= 100.0

		modulesProgress = append(modulesProgress, dto.ModuleProgressDetail{
			ModuleID:    module.ID,
			ModuleTitle: module.Title,
			Percentage:  modulePercentage,
			IsCompleted: isCompleted,
		})

		if isCompleted {
			completedModules++
		}
	}

	// Calculate overall course progress
	overallProgress := float64(completedModules) / float64(len(modules)) * 100.0

	// Create comprehensive response
	summary := &dto.CourseProgressSummary{
		CourseID:        course.ID,
		CourseTitle:     course.Title,
		TotalPercentage: overallProgress,
		IsCompleted:     overallProgress >= 100.0,
		ModulesProgress: modulesProgress,
	}

	return summary, nil
}

// Helper methods for optimized batch processing

func (s *userProgressService) batchGetContentsByModuleIDs(moduleIDs []uint) ([]*models.Content, error) {
	// This assumes the repository has a method to batch fetch contents by module IDs
	// If not available, we'll need to make individual queries but cache results
	var allContents []*models.Content
	for _, moduleID := range moduleIDs {
		contents, err := s.store.Contents.GetByModuleID(moduleID)
		if err != nil {
			continue // Skip modules with errors, don't fail the entire request
		}
		allContents = append(allContents, contents...)
	}
	return allContents, nil
}

func (s *userProgressService) batchGetEvaluationsByModuleIDs(moduleIDs []uint) ([]*models.Evaluation, error) {
	var allEvaluations []*models.Evaluation
	for _, moduleID := range moduleIDs {
		evaluations, err := s.store.Evaluations.GetByModuleID(moduleID)
		if err != nil {
			continue // Skip modules with errors
		}
		allEvaluations = append(allEvaluations, evaluations...)
	}
	return allEvaluations, nil
}

func (s *userProgressService) batchGetUserEvaluationAttempts(userID uint, moduleIDs []uint) ([]*models.EvaluationAttempt, error) {
	// Get all attempts for this user
	allAttempts, err := s.store.EvaluationAttempts.GetAll()
	if err != nil {
		return nil, err
	}

	// Filter attempts for the specific user and evaluations in the course modules
	// First, get all evaluation IDs for the modules
	evaluationIDs := make(map[uint]bool)
	for _, moduleID := range moduleIDs {
		evaluations, err := s.store.Evaluations.GetByModuleID(moduleID)
		if err != nil {
			continue
		}
		for _, eval := range evaluations {
			evaluationIDs[eval.ID] = true
		}
	}

	var userAttempts []*models.EvaluationAttempt
	for _, attempt := range allAttempts {
		if attempt.UserID == userID && evaluationIDs[attempt.EvaluationID] {
			userAttempts = append(userAttempts, attempt)
		}
	}
	return userAttempts, nil
}

func (s *userProgressService) groupContentsByModule(contents []*models.Content) map[uint][]*models.Content {
	contentsByModule := make(map[uint][]*models.Content)
	for _, content := range contents {
		contentsByModule[content.ModuleID] = append(contentsByModule[content.ModuleID], content)
	}
	return contentsByModule
}

func (s *userProgressService) groupEvaluationsByModule(evaluations []*models.Evaluation) map[uint][]*models.Evaluation {
	evaluationsByModule := make(map[uint][]*models.Evaluation)
	for _, evaluation := range evaluations {
		evaluationsByModule[evaluation.ModuleID] = append(evaluationsByModule[evaluation.ModuleID], evaluation)
	}
	return evaluationsByModule
}

func (s *userProgressService) groupProgressByContent(progressRecords []*models.UserProgress) map[uint]*models.UserProgress {
	progressByContent := make(map[uint]*models.UserProgress)
	for _, progress := range progressRecords {
		// Keep the most recent progress record for each content
		if existing, exists := progressByContent[progress.ContentID]; !exists || progress.UpdatedAt.After(existing.UpdatedAt) {
			progressByContent[progress.ContentID] = progress
		}
	}
	return progressByContent
}

func (s *userProgressService) groupAttemptsByEvaluation(attempts []*models.EvaluationAttempt) map[uint][]*models.EvaluationAttempt {
	attemptsByEvaluation := make(map[uint][]*models.EvaluationAttempt)
	for _, attempt := range attempts {
		attemptsByEvaluation[attempt.EvaluationID] = append(attemptsByEvaluation[attempt.EvaluationID], attempt)
	}
	return attemptsByEvaluation
}

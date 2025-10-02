package repositories

import (
	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type UserProgressRepository interface {
	Get(id uint) (*models.UserProgress, error)
	Create(userprogress *models.UserProgress) error
	Update(userprogress *models.UserProgress) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.UserProgress, error)
	GetAllWithPreloads() ([]*models.UserProgress, error)
	GetByUserAndCourse(userID, courseID uint) ([]*models.UserProgress, error)
	GetByUserAndModule(userID, moduleID uint) ([]*models.UserProgress, error)
	GetByUserAndContent(userID, contentID uint) (*models.UserProgress, error)
	CountCompletedByUserAndCourse(userID, courseID uint) (int64, error)
	CountCompletedByUserAndModule(userID, moduleID uint) (int64, error)
	BatchCreate(progressItems []*models.UserProgress) error
	GetCourseProgressSummary(userID, courseID uint) (*dto.CourseProgressSummary, error)
}

type userprogressRepository struct {
	*Repository
}

func NewUserProgressRepository(r *Repository) UserProgressRepository {
	return &userprogressRepository{
		Repository: r,
	}
}

func (r *userprogressRepository) Create(userprogress *models.UserProgress) error {
	return r.db.Create(userprogress).Error
}

func (r *userprogressRepository) Get(id uint) (*models.UserProgress, error) {
	var userprogress models.UserProgress
	if err := r.db.First(&userprogress, id).Error; err != nil {
		return nil, err
	}
	return &userprogress, nil
}

func (r *userprogressRepository) Update(userprogress *models.UserProgress) error {
	return r.db.Model(userprogress).Clauses(clause.Returning{}).Updates(userprogress).Error
}

func (r *userprogressRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.UserProgress{}).Where("id = ?", id).Updates(data).Error
}

func (r *userprogressRepository) Delete(id uint) error {
	var userprogress models.UserProgress
	userprogress.ID = id
	return r.db.Delete(&userprogress).Error
}

func (r *userprogressRepository) GetAll() ([]*models.UserProgress, error) {
	var userprogresss []*models.UserProgress
	if err := r.db.Find(&userprogresss).Error; err != nil {
		return nil, err
	}
	return userprogresss, nil
}

func (r *userprogressRepository) GetAllWithPreloads() ([]*models.UserProgress, error) {
	var userprogresss []*models.UserProgress
	if err := r.db.
		Preload("User").
		Preload("Course").
		Preload("Module").
		Preload("Content").
		Find(&userprogresss).Error; err != nil {
		return nil, err
	}
	return userprogresss, nil
}

func (r *userprogressRepository) GetByUserAndCourse(userID, courseID uint) ([]*models.UserProgress, error) {
	var userProgress []*models.UserProgress
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).Find(&userProgress).Error; err != nil {
		return nil, err
	}
	return userProgress, nil
}

func (r *userprogressRepository) GetByUserAndModule(userID, moduleID uint) ([]*models.UserProgress, error) {
	var userProgress []*models.UserProgress
	if err := r.db.Where("user_id = ? AND module_id = ?", userID, moduleID).Find(&userProgress).Error; err != nil {
		return nil, err
	}
	return userProgress, nil
}

func (r *userprogressRepository) GetByUserAndContent(userID, contentID uint) (*models.UserProgress, error) {
	var userProgress models.UserProgress
	if err := r.db.Where("user_id = ? AND content_id = ?", userID, contentID).First(&userProgress).Error; err != nil {
		return nil, err
	}
	return &userProgress, nil
}

func (r *userprogressRepository) CountCompletedByUserAndCourse(userID, courseID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND course_id = ? AND completed_at IS NOT NULL", userID, courseID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userprogressRepository) CountCompletedByUserAndModule(userID, moduleID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND module_id = ? AND completed_at IS NOT NULL", userID, moduleID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userprogressRepository) BatchCreate(progressItems []*models.UserProgress) error {
	if len(progressItems) == 0 {
		return nil
	}
	// Use batch insert for better performance
	return r.db.CreateInBatches(progressItems, 100).Error
}

func (r *userprogressRepository) GetCourseProgressSummary(userID, courseID uint) (*dto.CourseProgressSummary, error) {
	// Internal struct to capture raw SQL results
	type ProgressData struct {
		CourseID            uint    `db:"course_id"`
		CourseTitle         string  `db:"course_title"`
		ModuleID            uint    `db:"module_id"`
		ModuleTitle         string  `db:"module_title"`
		TotalItems          int     `db:"total_items"`
		CompletedContents   int     `db:"completed_contents"`
		PassedEvaluations   int     `db:"passed_evaluations"`
		CompletedItems      int     `db:"completed_items"`
	}

	var progressData []ProgressData
	
	// Complex query that calculates all progress data in one go
	query := `
	WITH module_stats AS (
		-- Get content and evaluation counts per module
		SELECT 
			m.id as module_id,
			m.title as module_title,
			m.course_id,
			COALESCE(content_count, 0) + COALESCE(evaluation_count, 0) as total_items,
			COALESCE(content_count, 0) as content_items,
			COALESCE(evaluation_count, 0) as evaluation_items
		FROM modules m
		LEFT JOIN (
			SELECT module_id, COUNT(*) as content_count
			FROM contents
			GROUP BY module_id
		) c ON c.module_id = m.id
		LEFT JOIN (
			SELECT module_id, COUNT(*) as evaluation_count
			FROM evaluations
			GROUP BY module_id
		) e ON e.module_id = m.id
		WHERE m.course_id = ?
	),
	user_content_progress AS (
		-- Get completed content for this user and course
		SELECT 
			up.module_id,
			COUNT(*) as completed_contents
		FROM user_progress up
		WHERE up.user_id = ? 
		AND up.course_id = ?
		AND up.completed_at IS NOT NULL
		GROUP BY up.module_id
	),
	user_evaluation_progress AS (
		-- Get passed evaluations for this user in this course
		SELECT 
			e.module_id,
			COUNT(*) as passed_evaluations
		FROM evaluation_attempts ea
		INNER JOIN evaluations e ON ea.evaluation_id = e.id
		WHERE ea.user_id = ?
		AND e.module_id IN (SELECT id FROM modules WHERE course_id = ?)
		AND ea.passed = true
		AND ea.submitted_at IS NOT NULL
		GROUP BY e.module_id
	)
	SELECT 
		c.id as course_id,
		c.title as course_title,
		ms.module_id,
		ms.module_title,
		ms.total_items,
		COALESCE(ucp.completed_contents, 0) as completed_contents,
		COALESCE(uep.passed_evaluations, 0) as passed_evaluations,
		COALESCE(ucp.completed_contents, 0) + COALESCE(uep.passed_evaluations, 0) as completed_items
	FROM courses c
	INNER JOIN module_stats ms ON ms.course_id = c.id
	LEFT JOIN user_content_progress ucp ON ucp.module_id = ms.module_id
	LEFT JOIN user_evaluation_progress uep ON uep.module_id = ms.module_id
	WHERE c.id = ?
	ORDER BY ms.module_id
	`

	if err := r.db.Raw(query, courseID, userID, courseID, userID, courseID, courseID).Scan(&progressData).Error; err != nil {
		return nil, err
	}

	if len(progressData) == 0 {
		// Course not found or has no modules - get course info directly
		var course models.Course
		if err := r.db.First(&course, courseID).Error; err != nil {
			return nil, err
		}
		
		return &dto.CourseProgressSummary{
			CourseID:        course.ID,
			CourseTitle:     course.Title,
			TotalPercentage: 100.0,
			IsCompleted:     true,
			ModulesProgress: []dto.ModuleProgressDetail{},
		}, nil
	}

	// Process the results to build the response
	var modulesProgress []dto.ModuleProgressDetail
	completedModules := 0
	totalModules := len(progressData)

	// Get course info from first row
	courseID = progressData[0].CourseID
	courseTitle := progressData[0].CourseTitle

	for _, data := range progressData {
		var modulePercentage float64
		if data.TotalItems > 0 {
			modulePercentage = float64(data.CompletedItems) / float64(data.TotalItems) * 100.0
		} else {
			modulePercentage = 100.0 // No content means 100% complete
		}
		
		isCompleted := modulePercentage >= 100.0
		if isCompleted {
			completedModules++
		}

		modulesProgress = append(modulesProgress, dto.ModuleProgressDetail{
			ModuleID:    data.ModuleID,
			ModuleTitle: data.ModuleTitle,
			Percentage:  modulePercentage,
			IsCompleted: isCompleted,
		})
	}

	// Calculate overall course progress
	var overallProgress float64
	if totalModules > 0 {
		overallProgress = float64(completedModules) / float64(totalModules) * 100.0
	} else {
		overallProgress = 100.0
	}

	// Create comprehensive response
	summary := &dto.CourseProgressSummary{
		CourseID:        courseID,
		CourseTitle:     courseTitle,
		TotalPercentage: overallProgress,
		IsCompleted:     overallProgress >= 100.0,
		ModulesProgress: modulesProgress,
	}

	return summary, nil
}

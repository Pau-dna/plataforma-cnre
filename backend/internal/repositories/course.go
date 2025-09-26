package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourseRepository interface {
	Get(id uint) (*models.Course, error)
	Create(course *models.Course) error
	Update(course *models.Course) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Course, error)
	IncrementStudentCount(courseID uint) error
	DecrementStudentCount(courseID uint) error
	IncrementModuleCount(courseID uint) error
	DecrementModuleCount(courseID uint) error
}

type courseRepository struct {
	*Repository
}

func NewCourseRepository(r *Repository) CourseRepository {
	return &courseRepository{
		Repository: r,
	}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) Get(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Model(course).Clauses(clause.Returning{}).Updates(course).Error
}

func (r *courseRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Course{}).Where("id = ?", id).Updates(data).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all modules for this course
		var modules []*models.Module
		if err := tx.Where("course_id = ?", id).Find(&modules).Error; err != nil {
			return err
		}
		
		// For each module, delete its contents and evaluations with cascade
		for _, module := range modules {
			// Get all evaluations for this module and cascade delete them
			var evaluations []*models.Evaluation
			if err := tx.Where("module_id = ?", module.ID).Find(&evaluations).Error; err != nil {
				return err
			}
			
			for _, evaluation := range evaluations {
				// Delete answers for all questions in this evaluation
				var questions []*models.Question
				if err := tx.Where("evaluation_id = ?", evaluation.ID).Find(&questions).Error; err != nil {
					return err
				}
				for _, question := range questions {
					if err := tx.Where("question_id = ?", question.ID).Delete(&models.Answer{}).Error; err != nil {
						return err
					}
				}
				
				// Delete questions for this evaluation
				if err := tx.Where("evaluation_id = ?", evaluation.ID).Delete(&models.Question{}).Error; err != nil {
					return err
				}
				
				// Delete evaluation attempts for this evaluation
				if err := tx.Where("evaluation_id = ?", evaluation.ID).Delete(&models.EvaluationAttempt{}).Error; err != nil {
					return err
				}
			}
			
			// Delete all evaluations for this module
			if err := tx.Where("module_id = ?", module.ID).Delete(&models.Evaluation{}).Error; err != nil {
				return err
			}
			
			// Delete all contents for this module
			if err := tx.Where("module_id = ?", module.ID).Delete(&models.Content{}).Error; err != nil {
				return err
			}
		}
		
		// Delete all modules for this course
		if err := tx.Where("course_id = ?", id).Delete(&models.Module{}).Error; err != nil {
			return err
		}
		
		// Delete user progress for this course
		if err := tx.Where("course_id = ?", id).Delete(&models.UserProgress{}).Error; err != nil {
			return err
		}
		
		// Delete enrollments for this course
		if err := tx.Where("course_id = ?", id).Delete(&models.Enrollment{}).Error; err != nil {
			return err
		}
		
		// Finally delete the course itself
		var course models.Course
		course.ID = id
		return tx.Delete(&course).Error
	})
}

func (r *courseRepository) GetAll() ([]*models.Course, error) {
	var courses []*models.Course
	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) IncrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("student_count", r.db.Raw("student_count + 1")).Error
}

func (r *courseRepository) DecrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("student_count", r.db.Raw("GREATEST(0, student_count - 1)")).Error
}

func (r *courseRepository) IncrementModuleCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("module_count", r.db.Raw("module_count + 1")).Error
}

func (r *courseRepository) DecrementModuleCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("module_count", r.db.Raw("GREATEST(0, module_count - 1)")).Error
}

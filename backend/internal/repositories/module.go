package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModuleRepository interface {
	Get(id uint) (*models.Module, error)
	Create(module *models.Module) error
	Update(module *models.Module) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Module, error)
	GetByCourseID(courseID uint) ([]*models.Module, error)
	GetWithContent(id uint) (*models.Module, error)
	GetMaxOrderByCourseID(courseID uint) (int, error)
}

type moduleRepository struct {
	*Repository
}

func NewModuleRepository(r *Repository) ModuleRepository {
	return &moduleRepository{
		Repository: r,
	}
}

func (r *moduleRepository) Create(module *models.Module) error {
	// Get the maximum order for this course
	var maxOrder int
	err := r.db.Model(&models.Module{}).
		Where("course_id = ?", module.CourseID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error
	if err != nil {
		return err
	}

	// Set the next order
	module.Order = maxOrder + 1

	// Create the module
	return r.db.Create(module).Error
}

func (r *moduleRepository) Get(id uint) (*models.Module, error) {
	var module models.Module
	if err := r.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepository) Update(module *models.Module) error {
	return r.db.Model(module).Clauses(clause.Returning{}).Updates(module).Error
}

func (r *moduleRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Module{}).Where("id = ?", id).Updates(data).Error
}

func (r *moduleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all evaluations for this module and cascade delete them
		var evaluations []*models.Evaluation
		if err := tx.Where("module_id = ?", id).Find(&evaluations).Error; err != nil {
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
		if err := tx.Where("module_id = ?", id).Delete(&models.Evaluation{}).Error; err != nil {
			return err
		}
		
		// Delete all contents for this module
		if err := tx.Where("module_id = ?", id).Delete(&models.Content{}).Error; err != nil {
			return err
		}
		
		// Delete user progress for this module
		if err := tx.Where("module_id = ?", id).Delete(&models.UserProgress{}).Error; err != nil {
			return err
		}
		
		// Finally delete the module itself
		var module models.Module
		module.ID = id
		return tx.Delete(&module).Error
	})
}

func (r *moduleRepository) GetAll() ([]*models.Module, error) {
	var modules []*models.Module
	if err := r.db.Preload("Contents").Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *moduleRepository) GetByCourseID(courseID uint) ([]*models.Module, error) {
	var modules []*models.Module
	if err := r.db.Preload("Contents").Where("course_id = ?", courseID).Order("\"order\" ASC").Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *moduleRepository) GetWithContent(id uint) (*models.Module, error) {
	var module models.Module
	if err := r.db.Preload("Contents").First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepository) GetMaxOrderByCourseID(courseID uint) (int, error) {
	var maxOrder int
	err := r.db.Model(&models.Module{}).
		Where("course_id = ?", courseID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

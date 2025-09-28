package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EvaluationRepository interface {
	Get(id uint) (*models.Evaluation, error)
	Create(evaluation *models.Evaluation) error
	Update(evaluation *models.Evaluation) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Evaluation, error)
	GetByModuleID(moduleID uint) ([]*models.Evaluation, error)
	GetWithQuestions(id uint) (*models.Evaluation, error)
}

type evaluationRepository struct {
	*Repository
}

func NewEvaluationRepository(r *Repository) EvaluationRepository {
	return &evaluationRepository{
		Repository: r,
	}
}

func (r *evaluationRepository) Create(evaluation *models.Evaluation) error {
	return r.db.Create(evaluation).Error
}

func (r *evaluationRepository) Get(id uint) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	if err := r.db.First(&evaluation, id).Error; err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (r *evaluationRepository) Update(evaluation *models.Evaluation) error {
	return r.db.Model(evaluation).Clauses(clause.Returning{}).Updates(evaluation).Error
}

func (r *evaluationRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Evaluation{}).Where("id = ?", id).Updates(data).Error
}

func (r *evaluationRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete all questions for this evaluation (answers will be deleted in cascade)
		if err := tx.Where(&models.Question{EvaluationID: id}).Delete(&models.Question{}).Error; err != nil {
			return err
		}

		// Delete all evaluation attempts for this evaluation
		if err := tx.Where(&models.EvaluationAttempt{EvaluationID: id}).Delete(&models.EvaluationAttempt{}).Error; err != nil {
			return err
		}

		// Finally delete the evaluation itself
		return tx.Delete(&models.Evaluation{}, id).Error
	})
}

func (r *evaluationRepository) GetAll() ([]*models.Evaluation, error) {
	var evaluations []*models.Evaluation
	if err := r.db.Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

func (r *evaluationRepository) GetByModuleID(moduleID uint) ([]*models.Evaluation, error) {
	var evaluations []*models.Evaluation
	if err := r.db.Where("module_id = ?", moduleID).Order("\"order\" ASC").Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

func (r *evaluationRepository) GetWithQuestions(id uint) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	if err := r.db.Preload("Questions").First(&evaluation, id).Error; err != nil {
		return nil, err
	}
	return &evaluation, nil
}

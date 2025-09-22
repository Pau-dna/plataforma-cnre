package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type EvaluationRepository interface {
	Get(id uint) (*models.Evaluation, error)
	Create(evaluation *models.Evaluation) error
	Update(evaluation *models.Evaluation) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Evaluation, error)
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
	var evaluation models.Evaluation
	evaluation.ID = id
	return r.db.Delete(&evaluation).Error
}

func (r *evaluationRepository) GetAll() ([]*models.Evaluation, error) {
	var evaluations []*models.Evaluation
	if err := r.db.Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

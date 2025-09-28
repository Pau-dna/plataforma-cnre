package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type EvaluationAttemptRepository interface {
	Get(id uint) (*models.EvaluationAttempt, error)
	Create(evaluationattempt *models.EvaluationAttempt) error
	Update(evaluationattempt *models.EvaluationAttempt) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.EvaluationAttempt, error)
	GetByUserAndEvaluation(userID, evaluationID uint) ([]*models.EvaluationAttempt, error)
	CountCompletedAttempts(userID, evaluationID uint) (int64, error)
	GetInProgressAttempt(userID, evaluationID uint) (*models.EvaluationAttempt, error)
}

type evaluationattemptRepository struct {
	*Repository
}

func NewEvaluationAttemptRepository(r *Repository) EvaluationAttemptRepository {
	return &evaluationattemptRepository{
		Repository: r,
	}
}

func (r *evaluationattemptRepository) Create(evaluationattempt *models.EvaluationAttempt) error {
	return r.db.Create(evaluationattempt).Error
}

func (r *evaluationattemptRepository) Get(id uint) (*models.EvaluationAttempt, error) {
	var evaluationattempt models.EvaluationAttempt
	if err := r.db.First(&evaluationattempt, id).Error; err != nil {
		return nil, err
	}
	return &evaluationattempt, nil
}

func (r *evaluationattemptRepository) Update(evaluationattempt *models.EvaluationAttempt) error {
	return r.db.Model(evaluationattempt).Clauses(clause.Returning{}).Updates(evaluationattempt).Error
}

func (r *evaluationattemptRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.EvaluationAttempt{}).Where("id = ?", id).Updates(data).Error
}

func (r *evaluationattemptRepository) Delete(id uint) error {
	var evaluationattempt models.EvaluationAttempt
	evaluationattempt.ID = id
	return r.db.Delete(&evaluationattempt).Error
}

func (r *evaluationattemptRepository) GetAll() ([]*models.EvaluationAttempt, error) {
	var evaluationattempts []*models.EvaluationAttempt
	if err := r.db.Find(&evaluationattempts).Error; err != nil {
		return nil, err
	}
	return evaluationattempts, nil
}

func (r *evaluationattemptRepository) GetByUserAndEvaluation(userID, evaluationID uint) ([]*models.EvaluationAttempt, error) {
	var attempts []*models.EvaluationAttempt
	if err := r.db.Where("user_id = ? AND evaluation_id = ?", userID, evaluationID).
		Order("created_at DESC").Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *evaluationattemptRepository) CountCompletedAttempts(userID, evaluationID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.EvaluationAttempt{}).
		Where("user_id = ? AND evaluation_id = ? AND submitted_at IS NOT NULL", userID, evaluationID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *evaluationattemptRepository) GetInProgressAttempt(userID, evaluationID uint) (*models.EvaluationAttempt, error) {
	var attempt models.EvaluationAttempt
	if err := r.db.Where("user_id = ? AND evaluation_id = ? AND submitted_at IS NULL", userID, evaluationID).
		First(&attempt).Error; err != nil {
		return nil, err
	}
	return &attempt, nil
}

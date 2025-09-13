
package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type QuestionRepository interface {
	Get(id uint) (*models.Question, error)
	Create(question *models.Question) error
	Update(question *models.Question) error
	Delete(id uint) error
}

type questionRepository struct {
	*Repository
}

func NewQuestionRepository(r *Repository) QuestionRepository {
	return &questionRepository{
		Repository: r,
	}
}

func (r *questionRepository) Create(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) Get(id uint) (*models.Question, error) {
	var question models.Question
	if err := r.db.First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) Update(question *models.Question) error {
	return r.db.Model(question).Clauses(clause.Returning{}).Updates(question).Error
}

func (r *questionRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Question{}).Where("id = ?", id).Updates(data).Error
}

func (r *questionRepository) Delete(id uint) error {
	var question models.Question
	question.ID = id
	return r.db.Delete(&question).Error
}

func (r *questionRepository) GetAll() ([]*models.Question, error) {
	var questions []*models.Question
	if err := r.db.Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

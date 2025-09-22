package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type AnswerRepository interface {
	Get(id uint) (*models.Answer, error)
	Create(answer *models.Answer) error
	Update(answer *models.Answer) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Answer, error)
	GetByQuestionID(questionID uint) ([]*models.Answer, error)
}

type answerRepository struct {
	*Repository
}

func NewAnswerRepository(r *Repository) AnswerRepository {
	return &answerRepository{
		Repository: r,
	}
}

func (r *answerRepository) Create(answer *models.Answer) error {
	return r.db.Create(answer).Error
}

func (r *answerRepository) Get(id uint) (*models.Answer, error) {
	var answer models.Answer
	if err := r.db.First(&answer, id).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *answerRepository) Update(answer *models.Answer) error {
	return r.db.Model(answer).Clauses(clause.Returning{}).Updates(answer).Error
}

func (r *answerRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Answer{}).Where("id = ?", id).Updates(data).Error
}

func (r *answerRepository) Delete(id uint) error {
	var answer models.Answer
	answer.ID = id
	return r.db.Delete(&answer).Error
}

func (r *answerRepository) GetAll() ([]*models.Answer, error) {
	var answers []*models.Answer
	if err := r.db.Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *answerRepository) GetByQuestionID(questionID uint) ([]*models.Answer, error) {
	var answers []*models.Answer
	if err := r.db.Where("question_id = ?", questionID).Order("\"order\" ASC").Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

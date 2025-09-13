
package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	Get(id uint) (*models.Content, error)
	Create(content *models.Content) error
	Update(content *models.Content) error
	Delete(id uint) error
}

type contentRepository struct {
	*Repository
}

func NewContentRepository(r *Repository) ContentRepository {
	return &contentRepository{
		Repository: r,
	}
}

func (r *contentRepository) Create(content *models.Content) error {
	return r.db.Create(content).Error
}

func (r *contentRepository) Get(id uint) (*models.Content, error) {
	var content models.Content
	if err := r.db.First(&content, id).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *contentRepository) Update(content *models.Content) error {
	return r.db.Model(content).Clauses(clause.Returning{}).Updates(content).Error
}

func (r *contentRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Content{}).Where("id = ?", id).Updates(data).Error
}

func (r *contentRepository) Delete(id uint) error {
	var content models.Content
	content.ID = id
	return r.db.Delete(&content).Error
}

func (r *contentRepository) GetAll() ([]*models.Content, error) {
	var contents []*models.Content
	if err := r.db.Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

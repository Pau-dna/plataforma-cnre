package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type ModuleRepository interface {
	Get(id uint) (*models.Module, error)
	Create(module *models.Module) error
	Update(module *models.Module) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Module, error)
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
	var module models.Module
	module.ID = id
	return r.db.Delete(&module).Error
}

func (r *moduleRepository) GetAll() ([]*models.Module, error) {
	var modules []*models.Module
	if err := r.db.Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

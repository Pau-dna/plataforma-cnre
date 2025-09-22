package repositories

import (
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
	GetByUserAndCourse(userID, courseID uint) ([]*models.UserProgress, error)
	GetByUserAndModule(userID, moduleID uint) ([]*models.UserProgress, error)
	GetByUserAndContent(userID, contentID uint) (*models.UserProgress, error)
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

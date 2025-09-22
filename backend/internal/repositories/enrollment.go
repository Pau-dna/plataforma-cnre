package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type EnrollmentRepository interface {
	Get(id uint) (*models.Enrollment, error)
	Create(enrollment *models.Enrollment) error
	Update(enrollment *models.Enrollment) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetUserEnrollment(userID uint, courseID uint) (*models.Enrollment, error)
	GetAll() ([]*models.Enrollment, error)
}

type enrollmentRepository struct {
	*Repository
}

func NewEnrollmentRepository(r *Repository) EnrollmentRepository {
	return &enrollmentRepository{
		Repository: r,
	}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) Get(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetUserEnrollment(userID uint, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Where(&models.Enrollment{UserID: userID, CourseID: courseID}).First(&enrollment).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Model(enrollment).Clauses(clause.Returning{}).Updates(enrollment).Error
}

func (r *enrollmentRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Enrollment{}).Where("id = ?", id).Updates(data).Error
}

func (r *enrollmentRepository) Delete(id uint) error {
	var enrollment models.Enrollment
	enrollment.ID = id
	return r.db.Delete(&enrollment).Error
}

func (r *enrollmentRepository) GetAll() ([]*models.Enrollment, error) {
	var enrollments []*models.Enrollment
	if err := r.db.Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

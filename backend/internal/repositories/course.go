package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type CourseRepository interface {
	Get(id uint) (*models.Course, error)
	Create(course *models.Course) error
	Update(course *models.Course) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Course, error)
	IncrementStudentCount(courseID uint) error
	DecrementStudentCount(courseID uint) error
	IncrementModuleCount(courseID uint) error
	DecrementModuleCount(courseID uint) error
}

type courseRepository struct {
	*Repository
}

func NewCourseRepository(r *Repository) CourseRepository {
	return &courseRepository{
		Repository: r,
	}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) Get(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Model(course).Clauses(clause.Returning{}).Updates(course).Error
}

func (r *courseRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Course{}).Where("id = ?", id).Updates(data).Error
}

func (r *courseRepository) Delete(id uint) error {
	var course models.Course
	course.ID = id
	return r.db.Delete(&course).Error
}

func (r *courseRepository) GetAll() ([]*models.Course, error) {
	var courses []*models.Course
	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) IncrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("student_count", r.db.Raw("student_count + 1")).Error
}

func (r *courseRepository) DecrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("student_count", r.db.Raw("GREATEST(0, student_count - 1)")).Error
}

func (r *courseRepository) IncrementModuleCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("module_count", r.db.Raw("module_count + 1")).Error
}

func (r *courseRepository) DecrementModuleCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("module_count", r.db.Raw("GREATEST(0, module_count - 1)")).Error
}

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
	GetAllWithCounts() ([]*models.Course, error)
	GetWithModules(id uint) (*models.Course, error)
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
	if err := r.db.Select("courses.*, "+
		"(SELECT COUNT(*) FROM modules WHERE modules.course_id = courses.id) as module_count, "+
		"(SELECT COUNT(*) FROM enrollments WHERE enrollments.course_id = courses.id) as student_count").
		First(&course, id).Error; err != nil {
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

func (r *courseRepository) GetAllWithCounts() ([]*models.Course, error) {
	var courses []*models.Course
	if err := r.db.Select("courses.*, " +
		"(SELECT COUNT(*) FROM modules WHERE modules.course_id = courses.id) as module_count, " +
		"(SELECT COUNT(*) FROM enrollments WHERE enrollments.course_id = courses.id) as student_count").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) GetWithModules(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.Preload("Modules").First(&course, id).Error; err != nil {
		return nil, err
	}

	// Calculate counts
	course.ModuleCount = len(course.Modules)

	// Get enrollment count
	var enrollmentCount int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ?", id).Count(&enrollmentCount).Error; err != nil {
		return nil, err
	}
	course.StudentCount = int(enrollmentCount)

	return &course, nil
}

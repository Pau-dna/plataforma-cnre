package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type EnrollmentRepository interface {
	Get(id uint) (*models.Enrollment, error)
	GetWithPreloads(id uint) (*models.Enrollment, error)
	Create(enrollment *models.Enrollment) error
	Update(enrollment *models.Enrollment) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetUserEnrollment(userID uint, courseID uint) (*models.Enrollment, error)
	GetAll() ([]*models.Enrollment, error)
	GetByUserID(userID uint) ([]*models.Enrollment, error)
	GetByCourseID(courseID uint) ([]*models.Enrollment, error)
	GetCourseKPIs(courseID uint) (int, float64, float64, string, error)
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
	if err := r.db.Preload("User").Preload("Course").First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetWithPreloads(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Preload("Course").Preload("User").First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetUserEnrollment(userID uint, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Preload("User").Preload("Course").Where(&models.Enrollment{UserID: userID, CourseID: courseID}).First(&enrollment).Error; err != nil {
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
	if err := r.db.Preload("User").Preload("Course").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetByUserID(userID uint) ([]*models.Enrollment, error) {
	var enrollments []*models.Enrollment
	if err := r.db.Preload("User").Preload("Course").Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetByCourseID(courseID uint) ([]*models.Enrollment, error) {
	var enrollments []*models.Enrollment
	if err := r.db.Preload("User").Preload("Course").Where("course_id = ?", courseID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetCourseKPIs(courseID uint) (int, float64, float64, string, error) {
	// Get course title
	var course models.Course
	if err := r.db.First(&course, courseID).Error; err != nil {
		return 0, 0, 0, "", err
	}
	
	// Count total enrollments for this course
	var totalEnrollments int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&totalEnrollments).Error; err != nil {
		return 0, 0, 0, course.Title, err
	}
	
	if totalEnrollments == 0 {
		return 0, 0, 0, course.Title, nil
	}
	
	// Count completed enrollments (progress = 100)
	var completedEnrollments int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ? AND progress = ?", courseID, 100.0).Count(&completedEnrollments).Error; err != nil {
		return int(totalEnrollments), 0, 0, course.Title, err
	}
	
	// Calculate average progress
	var avgProgress struct {
		AverageProgress float64 `db:"avg_progress"`
	}
	
	if err := r.db.Model(&models.Enrollment{}).
		Select("AVG(progress) as avg_progress").
		Where("course_id = ?", courseID).
		Scan(&avgProgress).Error; err != nil {
		return int(totalEnrollments), 0, 0, course.Title, err
	}
	
	// Calculate completion rate as percentage
	completionRate := float64(completedEnrollments) / float64(totalEnrollments) * 100
	
	return int(totalEnrollments), completionRate, avgProgress.AverageProgress, course.Title, nil
}

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
	"gorm.io/gorm"
)

type EnrollmentService interface {
	CreateEnrollment(userID, courseID uint) (*models.Enrollment, error)
	GetEnrollment(id uint) (*models.Enrollment, error)
	UpdateEnrollment(id uint, enrollment *models.Enrollment) (*models.Enrollment, error)
	UpdateEnrollmentPatch(id uint, data map[string]interface{}) (*models.Enrollment, error)
	DeleteEnrollment(id uint) error
	GetUserEnrollments(userID uint) ([]*models.Enrollment, error)
	GetCourseEnrollments(courseID uint) ([]*models.Enrollment, error)
	GetUserCourseEnrollment(userID, courseID uint) (*models.Enrollment, error)
	CompleteEnrollment(userID, courseID uint) error
	UpdateProgress(userID, courseID uint, progress float64) error
}

type enrollmentService struct {
	*Service
}

func NewEnrollmentService(service *Service) EnrollmentService {
	return &enrollmentService{
		Service: service,
	}
}

func (s *enrollmentService) CreateEnrollment(userID, courseID uint) (*models.Enrollment, error) {
	// Verify user exists
	_, err := s.store.Users.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Verify course exists
	_, err = s.store.Courses.Get(courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Check if enrollment already exists
	existing, _ := s.GetUserCourseEnrollment(userID, courseID)
	if existing != nil {
		return nil, fmt.Errorf("user is already enrolled in this course")
	}

	enrollment := &models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		EnrolledAt: time.Now(),
		Progress:   0.0,
	}

	// Use transaction to ensure consistency between enrollment creation and counter update
	var createdEnrollment *models.Enrollment
	err = s.store.DB().Transaction(func(tx *gorm.DB) error {
		// Create the enrollment
		if err := tx.Create(enrollment).Error; err != nil {
			return fmt.Errorf("failed to create enrollment: %w", err)
		}
		createdEnrollment = enrollment

		// Increment student count for the course
		if err := tx.Model(&models.Course{}).Where("id = ?", courseID).
			Update("student_count", tx.Raw("student_count + 1")).Error; err != nil {
			return fmt.Errorf("failed to increment student count: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdEnrollment, nil
}

func (s *enrollmentService) GetEnrollment(id uint) (*models.Enrollment, error) {
	enrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return nil, fmt.Errorf("enrollment not found: %w", err)
	}
	return enrollment, nil
}

func (s *enrollmentService) UpdateEnrollment(id uint, enrollmentData *models.Enrollment) (*models.Enrollment, error) {
	existingEnrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return nil, fmt.Errorf("enrollment not found: %w", err)
	}

	// Update fields
	existingEnrollment.Progress = enrollmentData.Progress
	existingEnrollment.CompletedAt = enrollmentData.CompletedAt

	if err := s.store.Enrollments.Update(existingEnrollment); err != nil {
		return nil, fmt.Errorf("failed to update enrollment: %w", err)
	}

	return existingEnrollment, nil
}

func (s *enrollmentService) UpdateEnrollmentPatch(enrollmentID uint, data map[string]interface{}) (*models.Enrollment, error) {
	if enrollmentID == 0 {
		return nil, errors.New("enrollment ID cannot be zero")
	}

	var enrollment dto.UpdateEnrollmentRequest
	if err := utils.MapToStructStrict(data, &enrollment); err != nil {
		return nil, errors.New("invalid data: " + err.Error())
	}

	if err := s.store.Enrollments.Patch(enrollmentID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Enrollments.Get(enrollmentID)
	if err != nil {
		return nil, errors.New("enrollment not found")
	}

	return updated, nil
}

func (s *enrollmentService) DeleteEnrollment(id uint) error {
	// Get the enrollment to get the course ID before deleting
	enrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get enrollment: %w", err)
	}

	courseID := enrollment.CourseID

	// Use transaction to ensure consistency between enrollment deletion and counter update
	err = s.store.DB().Transaction(func(tx *gorm.DB) error {
		// Delete the enrollment
		if err := tx.Delete(&models.Enrollment{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete enrollment: %w", err)
		}

		// Decrement student count for the course
		if err := tx.Model(&models.Course{}).Where("id = ?", courseID).
			Update("student_count", tx.Raw("GREATEST(0, student_count - 1)")).Error; err != nil {
			return fmt.Errorf("failed to decrement student count: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *enrollmentService) GetUserEnrollments(userID uint) ([]*models.Enrollment, error) {
	// Use the new repository method to filter by user ID at database level
	enrollments, err := s.store.Enrollments.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}

	return enrollments, nil
}

func (s *enrollmentService) GetCourseEnrollments(courseID uint) ([]*models.Enrollment, error) {
	// Use the new repository method to filter by course ID at database level
	enrollments, err := s.store.Enrollments.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}

	return enrollments, nil
}

func (s *enrollmentService) GetUserCourseEnrollment(userID, courseID uint) (*models.Enrollment, error) {
	// This would require a repository method to filter by both user ID and course ID
	enrollment, err := s.store.Enrollments.GetUserEnrollment(userID, courseID)
	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *enrollmentService) CompleteEnrollment(userID, courseID uint) error {
	enrollment, err := s.GetUserCourseEnrollment(userID, courseID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	now := time.Now()
	enrollment.CompletedAt = now
	enrollment.Progress = 100.0

	if err := s.store.Enrollments.Update(enrollment); err != nil {
		return fmt.Errorf("failed to complete enrollment: %w", err)
	}

	return nil
}

func (s *enrollmentService) UpdateProgress(userID, courseID uint, progress float64) error {
	enrollment, err := s.GetUserCourseEnrollment(userID, courseID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	// Validate progress value
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	enrollment.Progress = progress

	// If progress is 100%, mark as completed
	if progress >= 100.0 && enrollment.CompletedAt.IsZero() {
		enrollment.CompletedAt = time.Now()
	}

	if err := s.store.Enrollments.Update(enrollment); err != nil {
		return fmt.Errorf("failed to update enrollment progress: %w", err)
	}

	return nil
}

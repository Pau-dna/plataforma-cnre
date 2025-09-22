package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
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

	if err := s.store.Enrollments.Create(enrollment); err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	return enrollment, nil
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
	if err := s.store.Enrollments.Delete(id); err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	return nil
}

func (s *enrollmentService) GetUserEnrollments(userID uint) ([]*models.Enrollment, error) {
	// This would require a repository method to filter by user ID
	// For now, we'll implement a basic version
	enrollments, err := s.store.Enrollments.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}

	// Filter by user ID
	var userEnrollments []*models.Enrollment
	for _, enrollment := range enrollments {
		if enrollment.UserID == userID {
			userEnrollments = append(userEnrollments, enrollment)
		}
	}

	return userEnrollments, nil
}

func (s *enrollmentService) GetCourseEnrollments(courseID uint) ([]*models.Enrollment, error) {
	// This would require a repository method to filter by course ID
	// For now, we'll implement a basic version
	enrollments, err := s.store.Enrollments.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}

	// Filter by course ID
	var courseEnrollments []*models.Enrollment
	for _, enrollment := range enrollments {
		if enrollment.CourseID == courseID {
			courseEnrollments = append(courseEnrollments, enrollment)
		}
	}

	return courseEnrollments, nil
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

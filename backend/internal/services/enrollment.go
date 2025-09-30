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
	GetEnrollmentWithPreloads(id uint) (*models.Enrollment, error)
	UpdateEnrollment(id uint, enrollment *models.Enrollment) (*models.Enrollment, error)
	UpdateEnrollmentPatch(id uint, data map[string]interface{}) (*models.Enrollment, error)
	DeleteEnrollment(id uint) error
	GetUserEnrollments(userID uint) ([]*models.Enrollment, error)
	GetCourseEnrollments(courseID uint) ([]*models.Enrollment, error)
	GetUserCourseEnrollment(userID, courseID uint) (*models.Enrollment, error)
	CompleteEnrollment(userID, courseID uint) error
	UpdateProgress(userID, courseID uint, progress float64) error
	GetCourseKPIs(courseID uint) (*dto.CourseKPIResponse, error)
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
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verify course exists
	_, err = s.store.Courses.Get(courseID)
	if err != nil {
		return nil, fmt.Errorf("curso no encontrado: %w", err)
	}

	// Check if enrollment already exists
	existing, _ := s.GetUserCourseEnrollment(userID, courseID)
	if existing != nil {
		return nil, fmt.Errorf("el usuario ya está inscrito en este curso")
	}

	enrollment := &models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		EnrolledAt: time.Now(),
		Progress:   0.0,
	}

	if err := s.store.Enrollments.Create(enrollment); err != nil {
		return nil, fmt.Errorf("error al crear la inscripción: %w", err)
	}

	// Increment student count for the course
	if err := s.store.Courses.IncrementStudentCount(courseID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to increment student count for course %d: %v\n", courseID, err)
	}

	return enrollment, nil
}

func (s *enrollmentService) GetEnrollment(id uint) (*models.Enrollment, error) {
	enrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return nil, fmt.Errorf("inscripción no encontrada: %w", err)
	}
	return enrollment, nil
}

func (s *enrollmentService) GetEnrollmentWithPreloads(id uint) (*models.Enrollment, error) {
	enrollment, err := s.store.Enrollments.GetWithPreloads(id)
	if err != nil {
		return nil, fmt.Errorf("inscripción no encontrada: %w", err)
	}
	return enrollment, nil
}

func (s *enrollmentService) UpdateEnrollment(id uint, enrollmentData *models.Enrollment) (*models.Enrollment, error) {
	existingEnrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return nil, fmt.Errorf("inscripción no encontrada: %w", err)
	}

	// Update fields
	existingEnrollment.Progress = enrollmentData.Progress
	existingEnrollment.CompletedAt = enrollmentData.CompletedAt

	if err := s.store.Enrollments.Update(existingEnrollment); err != nil {
		return nil, fmt.Errorf("error al actualizar la inscripción: %w", err)
	}

	return existingEnrollment, nil
}

func (s *enrollmentService) UpdateEnrollmentPatch(enrollmentID uint, data map[string]interface{}) (*models.Enrollment, error) {
	if enrollmentID == 0 {
		return nil, errors.New("el ID de inscripción no puede ser cero")
	}

	var enrollment dto.UpdateEnrollmentRequest
	if err := utils.MapToStructStrict(data, &enrollment); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.Enrollments.Patch(enrollmentID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Enrollments.Get(enrollmentID)
	if err != nil {
		return nil, errors.New("inscripción no encontrada")
	}

	return updated, nil
}

func (s *enrollmentService) DeleteEnrollment(id uint) error {
	// Get the enrollment to get the course ID before deleting
	enrollment, err := s.store.Enrollments.Get(id)
	if err != nil {
		return fmt.Errorf("error al obtener la inscripción: %w", err)
	}

	courseID := enrollment.CourseID

	if err := s.store.Enrollments.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar la inscripción: %w", err)
	}

	// Decrement student count for the course
	if err := s.store.Courses.DecrementStudentCount(courseID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to decrement student count for course %d: %v\n", courseID, err)
	}

	return nil
}

func (s *enrollmentService) GetUserEnrollments(userID uint) ([]*models.Enrollment, error) {
	// Use the new repository method to filter by user ID at database level
	enrollments, err := s.store.Enrollments.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las inscripciones: %w", err)
	}

	return enrollments, nil
}

func (s *enrollmentService) GetCourseEnrollments(courseID uint) ([]*models.Enrollment, error) {
	// Use the new repository method to filter by course ID at database level
	enrollments, err := s.store.Enrollments.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las inscripciones: %w", err)
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
		return fmt.Errorf("inscripción no encontrada: %w", err)
	}

	now := time.Now()
	enrollment.CompletedAt = now
	enrollment.Progress = 100.0

	if err := s.store.Enrollments.Update(enrollment); err != nil {
		return fmt.Errorf("error al completar la inscripción: %w", err)
	}

	return nil
}

func (s *enrollmentService) UpdateProgress(userID, courseID uint, progress float64) error {
	enrollment, err := s.GetUserCourseEnrollment(userID, courseID)
	if err != nil {
		return fmt.Errorf("inscripción no encontrada: %w", err)
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
		return fmt.Errorf("error al actualizar el progreso de la inscripción: %w", err)
	}

	return nil
}

func (s *enrollmentService) GetCourseKPIs(courseID uint) (*dto.CourseKPIResponse, error) {
	studentCount, completionRate, avgProgress, courseTitle, err := s.store.Enrollments.GetCourseKPIs(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener KPIs del curso: %w", err)
	}

	return &dto.CourseKPIResponse{
		CourseID:        courseID,
		CourseTitle:     courseTitle,
		StudentCount:    studentCount,
		CompletionRate:  completionRate,
		AverageProgress: avgProgress,
	}, nil
}

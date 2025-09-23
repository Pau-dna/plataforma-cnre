package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type CourseService interface {
	CreateCourse(course *models.Course) (*models.Course, error)
	GetCourse(id uint) (*models.Course, error)
	UpdateCourse(id uint, course *models.Course) (*models.Course, error)
	UpdateCoursePatch(id uint, data map[string]interface{}) (*models.Course, error)
	DeleteCourse(id uint) error
	GetAllCourses() ([]*models.Course, error)
	GetCourseWithModules(id uint) (*models.Course, error)
	GetCoursesWithEnrollmentCount() ([]*models.Course, error)
}

type courseService struct {
	*Service
}

func NewCourseService(service *Service) CourseService {
	return &courseService{
		Service: service,
	}
}

func (s *courseService) CreateCourse(course *models.Course) (*models.Course, error) {
	if err := s.store.Courses.Create(course); err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}
	return course, nil
}

func (s *courseService) GetCourse(id uint) (*models.Course, error) {
	course, err := s.store.Courses.Get(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	return course, nil
}

func (s *courseService) UpdateCourse(id uint, courseData *models.Course) (*models.Course, error) {
	existingCourse, err := s.store.Courses.Get(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Update fields
	existingCourse.Title = courseData.Title
	existingCourse.Description = courseData.Description
	existingCourse.ShortDescription = courseData.ShortDescription
	existingCourse.ImageURL = courseData.ImageURL

	if err := s.store.Courses.Update(existingCourse); err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	return existingCourse, nil
}

func (s *courseService) UpdateCoursePatch(courseID uint, data map[string]interface{}) (*models.Course, error) {
	if courseID == 0 {
		return nil, errors.New("course ID cannot be zero")
	}

	var course dto.UpdateCourseRequest
	if err := utils.MapToStructStrict(data, &course); err != nil {
		return nil, errors.New("invalid data: " + err.Error())
	}

	if err := s.store.Courses.Patch(courseID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Courses.Get(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	return updated, nil
}

func (s *courseService) DeleteCourse(id uint) error {
	if err := s.store.Courses.Delete(id); err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}

func (s *courseService) GetAllCourses() ([]*models.Course, error) {
	courses, err := s.store.Courses.GetAllWithCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

func (s *courseService) GetCourseWithModules(id uint) (*models.Course, error) {
	course, err := s.store.Courses.GetWithModules(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	return course, nil
}

func (s *courseService) GetCoursesWithEnrollmentCount() ([]*models.Course, error) {
	courses, err := s.store.Courses.GetAllWithCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

package services

import (
	"fmt"

	"github.com/imlargo/go-api-template/internal/models"
)

type CourseService interface {
	CreateCourse(course *models.Course) (*models.Course, error)
	GetCourse(id uint) (*models.Course, error)
	UpdateCourse(id uint, course *models.Course) (*models.Course, error)
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

func (s *courseService) DeleteCourse(id uint) error {
	if err := s.store.Courses.Delete(id); err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}

func (s *courseService) GetAllCourses() ([]*models.Course, error) {
	courses, err := s.store.Courses.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

func (s *courseService) GetCourseWithModules(id uint) (*models.Course, error) {
	// This would require a repository method to preload modules
	course, err := s.store.Courses.Get(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	
	// For now, return the course - would need to implement preloading in repository
	return course, nil
}

func (s *courseService) GetCoursesWithEnrollmentCount() ([]*models.Course, error) {
	// This would require a more complex repository method to calculate enrollment counts
	courses, err := s.store.Courses.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}
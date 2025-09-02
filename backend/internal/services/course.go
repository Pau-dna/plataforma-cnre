package services

import (
	"github.com/imlargo/go-api-template/internal/models"
)

type CourseService interface {
}

type courseService struct {
	*Service
}

func NewCourseService(
	service *Service,
) CourseService {
	return &courseService{
		Service: service,
	}
}

func (s *courseService) GetCourse(id uint) (*models.Course, error) {
	return nil, nil
}

package services

import (
	"fmt"

	"github.com/imlargo/go-api-template/internal/models"
)

type ModuleService interface {
	CreateModule(module *models.Module) (*models.Module, error)
	GetModule(id uint) (*models.Module, error)
	UpdateModule(id uint, module *models.Module) (*models.Module, error)
	DeleteModule(id uint) error
	GetModulesByCourse(courseID uint) ([]*models.Module, error)
	GetModuleWithContent(id uint) (*models.Module, error)
	ReorderModules(courseID uint, moduleOrders []struct{ ID uint; Order int }) error
}

type moduleService struct {
	*Service
}

func NewModuleService(service *Service) ModuleService {
	return &moduleService{
		Service: service,
	}
}

func (s *moduleService) CreateModule(module *models.Module) (*models.Module, error) {
	// Verify course exists
	_, err := s.store.Courses.Get(module.CourseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	if err := s.store.Modules.Create(module); err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}
	return module, nil
}

func (s *moduleService) GetModule(id uint) (*models.Module, error) {
	module, err := s.store.Modules.Get(id)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}
	return module, nil
}

func (s *moduleService) UpdateModule(id uint, moduleData *models.Module) (*models.Module, error) {
	existingModule, err := s.store.Modules.Get(id)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	// Update fields
	existingModule.Title = moduleData.Title
	existingModule.Description = moduleData.Description
	existingModule.Order = moduleData.Order

	if err := s.store.Modules.Update(existingModule); err != nil {
		return nil, fmt.Errorf("failed to update module: %w", err)
	}

	return existingModule, nil
}

func (s *moduleService) DeleteModule(id uint) error {
	if err := s.store.Modules.Delete(id); err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}
	return nil
}

func (s *moduleService) GetModulesByCourse(courseID uint) ([]*models.Module, error) {
	// This would require a repository method to filter by course ID
	// For now, we'll implement a basic version
	modules, err := s.store.Modules.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	// Filter by course ID
	var courseModules []*models.Module
	for _, module := range modules {
		if module.CourseID == courseID {
			courseModules = append(courseModules, module)
		}
	}

	return courseModules, nil
}

func (s *moduleService) GetModuleWithContent(id uint) (*models.Module, error) {
	// This would require a repository method to preload content
	module, err := s.store.Modules.Get(id)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}
	
	// For now, return the module - would need to implement preloading in repository
	return module, nil
}

func (s *moduleService) ReorderModules(courseID uint, moduleOrders []struct{ ID uint; Order int }) error {
	// Verify course exists
	_, err := s.store.Courses.Get(courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// Update each module's order
	for _, order := range moduleOrders {
		module, err := s.store.Modules.Get(order.ID)
		if err != nil {
			continue // Skip invalid modules
		}
		
		if module.CourseID != courseID {
			continue // Skip modules from other courses
		}

		module.Order = order.Order
		if err := s.store.Modules.Update(module); err != nil {
			return fmt.Errorf("failed to update module order: %w", err)
		}
	}

	return nil
}
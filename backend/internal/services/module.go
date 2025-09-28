package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type ModuleService interface {
	CreateModule(module *models.Module) (*models.Module, error)
	GetModule(id uint) (*models.Module, error)
	UpdateModule(id uint, module *models.Module) (*models.Module, error)
	UpdateModulePatch(id uint, data map[string]interface{}) (*models.Module, error)
	DeleteModule(id uint) error
	GetModulesByCourse(courseID uint) ([]*models.Module, error)
	GetModuleWithContent(id uint) (*models.Module, error)
	ReorderModules(courseID uint, moduleOrders []struct {
		ID    uint
		Order int
	}) error
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
		return nil, fmt.Errorf("curso no encontrado: %w", err)
	}

	if err := s.store.Modules.Create(module); err != nil {
		return nil, fmt.Errorf("error al crear el módulo: %w", err)
	}

	// Increment module count for the course
	if err := s.store.Courses.IncrementModuleCount(module.CourseID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to increment module count for course %d: %v\n", module.CourseID, err)
	}

	return module, nil
}

func (s *moduleService) GetModule(id uint) (*models.Module, error) {
	module, err := s.store.Modules.Get(id)
	if err != nil {
		return nil, fmt.Errorf("módulo no encontrado: %w", err)
	}
	return module, nil
}

func (s *moduleService) UpdateModule(id uint, moduleData *models.Module) (*models.Module, error) {
	existingModule, err := s.store.Modules.Get(id)
	if err != nil {
		return nil, fmt.Errorf("módulo no encontrado: %w", err)
	}

	// Update fields
	existingModule.Title = moduleData.Title
	existingModule.Description = moduleData.Description
	existingModule.Order = moduleData.Order

	if err := s.store.Modules.Update(existingModule); err != nil {
		return nil, fmt.Errorf("error al actualizar el módulo: %w", err)
	}

	return existingModule, nil
}

func (s *moduleService) UpdateModulePatch(moduleID uint, data map[string]interface{}) (*models.Module, error) {
	if moduleID == 0 {
		return nil, errors.New("module ID cannot be zero")
	}

	var module dto.UpdateModuleRequest
	if err := utils.MapToStructStrict(data, &module); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.Modules.Patch(moduleID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Modules.Get(moduleID)
	if err != nil {
		return nil, errors.New("módulo no encontrado")
	}

	return updated, nil
}

func (s *moduleService) DeleteModule(id uint) error {
	// Get the module to get the course ID before deleting
	module, err := s.store.Modules.Get(id)
	if err != nil {
		return fmt.Errorf("error al obtener el módulo: %w", err)
	}

	courseID := module.CourseID

	if err := s.store.Modules.Delete(id); err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}

	// Decrement module count for the course
	if err := s.store.Courses.DecrementModuleCount(courseID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to decrement module count for course %d: %v\n", courseID, err)
	}

	return nil
}

func (s *moduleService) GetModulesByCourse(courseID uint) ([]*models.Module, error) {
	// Use the optimized repository method to filter by course ID at database level
	modules, err := s.store.Modules.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	return modules, nil
}

func (s *moduleService) GetModuleWithContent(id uint) (*models.Module, error) {
	// Use the new repository method to preload content
	module, err := s.store.Modules.GetWithContent(id)
	if err != nil {
		return nil, fmt.Errorf("módulo no encontrado: %w", err)
	}

	return module, nil
}

func (s *moduleService) ReorderModules(courseID uint, moduleOrders []struct {
	ID    uint
	Order int
}) error {
	// Verify course exists
	_, err := s.store.Courses.Get(courseID)
	if err != nil {
		return fmt.Errorf("curso no encontrado: %w", err)
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

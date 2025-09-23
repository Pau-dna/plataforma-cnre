package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
	"gorm.io/gorm"
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
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Use transaction to ensure consistency between module creation and counter update
	var createdModule *models.Module
	err = s.store.DB().Transaction(func(tx *gorm.DB) error {
		// Get the maximum order for this course
		var maxOrder int
		err := tx.Model(&models.Module{}).
			Where("course_id = ?", module.CourseID).
			Select("COALESCE(MAX(\"order\"), 0)").
			Scan(&maxOrder).Error
		if err != nil {
			return err
		}

		// Set the next order
		module.Order = maxOrder + 1

		// Create the module
		if err := tx.Create(module).Error; err != nil {
			return fmt.Errorf("failed to create module: %w", err)
		}
		createdModule = module

		// Increment module count for the course
		if err := tx.Model(&models.Course{}).Where("id = ?", module.CourseID).
			Update("module_count", tx.Raw("module_count + 1")).Error; err != nil {
			return fmt.Errorf("failed to increment module count: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdModule, nil
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

func (s *moduleService) UpdateModulePatch(moduleID uint, data map[string]interface{}) (*models.Module, error) {
	if moduleID == 0 {
		return nil, errors.New("module ID cannot be zero")
	}

	var module dto.UpdateModuleRequest
	if err := utils.MapToStructStrict(data, &module); err != nil {
		return nil, errors.New("invalid data: " + err.Error())
	}

	if err := s.store.Modules.Patch(moduleID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Modules.Get(moduleID)
	if err != nil {
		return nil, errors.New("module not found")
	}

	return updated, nil
}

func (s *moduleService) DeleteModule(id uint) error {
	// Get the module to get the course ID before deleting
	module, err := s.store.Modules.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get module: %w", err)
	}

	courseID := module.CourseID

	// Use transaction to ensure consistency between module deletion and counter update
	err = s.store.DB().Transaction(func(tx *gorm.DB) error {
		// Delete the module
		if err := tx.Delete(&models.Module{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete module: %w", err)
		}

		// Decrement module count for the course
		if err := tx.Model(&models.Course{}).Where("id = ?", courseID).
			Update("module_count", tx.Raw("GREATEST(0, module_count - 1)")).Error; err != nil {
			return fmt.Errorf("failed to decrement module count: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *moduleService) GetModulesByCourse(courseID uint) ([]*models.Module, error) {
	// Use the new repository method to filter by course ID at database level
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
		return nil, fmt.Errorf("module not found: %w", err)
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

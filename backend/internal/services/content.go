package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type ContentService interface {
	CreateContent(content *models.Content) (*models.Content, error)
	GetContent(id uint) (*models.Content, error)
	UpdateContent(id uint, content *models.Content) (*models.Content, error)
	UpdateContentPatch(id uint, data map[string]interface{}) (*models.Content, error)
	DeleteContent(id uint) error
	GetContentsByModule(moduleID uint) ([]*models.Content, error)
	ReorderContent(moduleID uint, contentOrders []struct {
		ID    uint
		Order int
	}) error
}

type contentService struct {
	*Service
}

func NewContentService(service *Service) ContentService {
	return &contentService{
		Service: service,
	}
}

func (s *contentService) CreateContent(content *models.Content) (*models.Content, error) {
	// Verify module exists
	_, err := s.store.Modules.Get(content.ModuleID)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	if err := s.store.Contents.Create(content); err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}
	return content, nil
}

func (s *contentService) GetContent(id uint) (*models.Content, error) {
	content, err := s.store.Contents.Get(id)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}
	return content, nil
}

func (s *contentService) UpdateContent(id uint, contentData *models.Content) (*models.Content, error) {
	existingContent, err := s.store.Contents.Get(id)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	// Update fields
	existingContent.Title = contentData.Title
	existingContent.Description = contentData.Description
	existingContent.Order = contentData.Order
	existingContent.Body = contentData.Body
	existingContent.MediaURL = contentData.MediaURL
	existingContent.Type = contentData.Type

	if err := s.store.Contents.Update(existingContent); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return existingContent, nil
}

func (s *contentService) UpdateContentPatch(contentID uint, data map[string]interface{}) (*models.Content, error) {
	if contentID == 0 {
		return nil, errors.New("content ID cannot be zero")
	}

	var content dto.UpdateContentRequest
	if err := utils.MapToStructStrict(data, &content); err != nil {
		return nil, errors.New("invalid data: " + err.Error())
	}

	if err := s.store.Contents.Patch(contentID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Contents.Get(contentID)
	if err != nil {
		return nil, errors.New("content not found")
	}

	return updated, nil
}

func (s *contentService) DeleteContent(id uint) error {
	if err := s.store.Contents.Delete(id); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}
	return nil
}

func (s *contentService) GetContentsByModule(moduleID uint) ([]*models.Content, error) {
	// Use the new repository method to filter by module ID at database level
	contents, err := s.store.Contents.GetByModuleID(moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}

	return contents, nil
}

func (s *contentService) ReorderContent(moduleID uint, contentOrders []struct {
	ID    uint
	Order int
}) error {
	// Verify module exists
	_, err := s.store.Modules.Get(moduleID)
	if err != nil {
		return fmt.Errorf("module not found: %w", err)
	}

	// Update each content's order
	for _, order := range contentOrders {
		content, err := s.store.Contents.Get(order.ID)
		if err != nil {
			continue // Skip invalid content
		}

		if content.ModuleID != moduleID {
			continue // Skip content from other modules
		}

		content.Order = order.Order
		if err := s.store.Contents.Update(content); err != nil {
			return fmt.Errorf("failed to update content order: %w", err)
		}
	}

	return nil
}

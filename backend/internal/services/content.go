package services

import (
	"fmt"

	"github.com/imlargo/go-api-template/internal/models"
)

type ContentService interface {
	CreateContent(content *models.Content) (*models.Content, error)
	GetContent(id uint) (*models.Content, error)
	UpdateContent(id uint, content *models.Content) (*models.Content, error)
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

func (s *contentService) DeleteContent(id uint) error {
	if err := s.store.Contents.Delete(id); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}
	return nil
}

func (s *contentService) GetContentsByModule(moduleID uint) ([]*models.Content, error) {
	// This would require a repository method to filter by module ID
	// For now, we'll implement a basic version
	contents, err := s.store.Contents.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}

	// Filter by module ID
	var moduleContents []*models.Content
	for _, content := range contents {
		if content.ModuleID == moduleID {
			moduleContents = append(moduleContents, content)
		}
	}

	return moduleContents, nil
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

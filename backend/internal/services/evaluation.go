package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type EvaluationService interface {
	CreateEvaluation(evaluation *models.Evaluation) (*models.Evaluation, error)
	GetEvaluation(id uint) (*models.Evaluation, error)
	UpdateEvaluation(id uint, evaluation *models.Evaluation) (*models.Evaluation, error)
	UpdateEvaluationPatch(id uint, data map[string]interface{}) (*models.Evaluation, error)
	DeleteEvaluation(id uint) error
	GetEvaluationsByModule(moduleID uint) ([]*models.Evaluation, error)
	GetEvaluationWithQuestions(id uint) (*models.Evaluation, error)
}

type evaluationService struct {
	*Service
}

func NewEvaluationService(service *Service) EvaluationService {
	return &evaluationService{
		Service: service,
	}
}

func (s *evaluationService) CreateEvaluation(evaluation *models.Evaluation) (*models.Evaluation, error) {
	// Verify module exists
	_, err := s.store.Modules.Get(evaluation.ModuleID)
	if err != nil {
		return nil, fmt.Errorf("módulo no encontrado: %w", err)
	}

	if err := s.store.Evaluations.Create(evaluation); err != nil {
		return nil, fmt.Errorf("error al crear la evaluación: %w", err)
	}
	return evaluation, nil
}

func (s *evaluationService) GetEvaluation(id uint) (*models.Evaluation, error) {
	evaluation, err := s.store.Evaluations.Get(id)
	if err != nil {
		return nil, fmt.Errorf("evaluación no encontrada: %w", err)
	}
	return evaluation, nil
}

func (s *evaluationService) UpdateEvaluation(id uint, evaluationData *models.Evaluation) (*models.Evaluation, error) {
	existingEvaluation, err := s.store.Evaluations.Get(id)
	if err != nil {
		return nil, fmt.Errorf("evaluación no encontrada: %w", err)
	}

	// Update fields
	existingEvaluation.Title = evaluationData.Title
	existingEvaluation.Description = evaluationData.Description
	existingEvaluation.Order = evaluationData.Order
	existingEvaluation.QuestionCount = evaluationData.QuestionCount
	existingEvaluation.PassingScore = evaluationData.PassingScore
	existingEvaluation.MaxAttempts = evaluationData.MaxAttempts
	existingEvaluation.TimeLimit = evaluationData.TimeLimit
	existingEvaluation.Type = evaluationData.Type

	if err := s.store.Evaluations.Update(existingEvaluation); err != nil {
		return nil, fmt.Errorf("error al actualizar la evaluación: %w", err)
	}

	return existingEvaluation, nil
}

func (s *evaluationService) UpdateEvaluationPatch(evaluationID uint, data map[string]interface{}) (*models.Evaluation, error) {
	if evaluationID == 0 {
		return nil, errors.New("evaluation ID cannot be zero")
	}

	var evaluation dto.UpdateEvaluationRequest
	if err := utils.MapToStructStrict(data, &evaluation); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.Evaluations.Patch(evaluationID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Evaluations.Get(evaluationID)
	if err != nil {
		return nil, errors.New("evaluación no encontrada")
	}

	return updated, nil
}

func (s *evaluationService) DeleteEvaluation(id uint) error {
	if err := s.store.Evaluations.Delete(id); err != nil {
		return fmt.Errorf("failed to delete evaluation: %w", err)
	}
	return nil
}

func (s *evaluationService) GetEvaluationsByModule(moduleID uint) ([]*models.Evaluation, error) {
	// Use the optimized repository method to filter by module ID at database level
	evaluations, err := s.store.Evaluations.GetByModuleID(moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get evaluations: %w", err)
	}

	return evaluations, nil
}

func (s *evaluationService) GetEvaluationWithQuestions(id uint) (*models.Evaluation, error) {
	// Use the new repository method to preload questions
	evaluation, err := s.store.Evaluations.GetWithQuestions(id)
	if err != nil {
		return nil, fmt.Errorf("evaluación no encontrada: %w", err)
	}

	return evaluation, nil
}

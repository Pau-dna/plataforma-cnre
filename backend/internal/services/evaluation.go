package services

import (
	"github.com/imlargo/go-api-template/internal/models"
)

type EvaluationService interface {
}

type evaluationService struct {
	*Service
}

func NewEvaluationService(
	service *Service,
) EvaluationService {
	return &evaluationService{
		Service: service,
	}
}

func (s *evaluationService) GetEvaluation(id uint) (*models.Evaluation, error) {
	return nil, nil
}

func (s *evaluationService) GenerateQuestions(id uint) (*[]models.Question, error) {
	return nil, nil
}

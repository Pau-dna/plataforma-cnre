package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type QuestionService interface {
	CreateQuestion(question *models.Question) (*models.Question, error)
	GetQuestion(id uint) (*models.Question, error)
	UpdateQuestion(id uint, question *models.Question) (*models.Question, error)
	UpdateQuestionPatch(id uint, data map[string]interface{}) (*models.Question, error)
	DeleteQuestion(id uint) error
	GetQuestionsByEvaluation(evaluationID uint) ([]*models.Question, error)
	GetQuestionWithAnswers(id uint) (*models.Question, error)
}

type questionService struct {
	*Service
}

func NewQuestionService(service *Service) QuestionService {
	return &questionService{
		Service: service,
	}
}

func (s *questionService) CreateQuestion(question *models.Question) (*models.Question, error) {
	// Verify evaluation exists
	_, err := s.store.Evaluations.Get(question.EvaluationID)
	if err != nil {
		return nil, fmt.Errorf("evaluación no encontrada: %w", err)
	}

	if err := s.store.Questions.Create(question); err != nil {
		return nil, fmt.Errorf("error al crear la pregunta: %w", err)
	}
	return question, nil
}

func (s *questionService) GetQuestion(id uint) (*models.Question, error) {
	question, err := s.store.Questions.Get(id)
	if err != nil {
		return nil, fmt.Errorf("pregunta no encontrada: %w", err)
	}
	return question, nil
}

func (s *questionService) UpdateQuestion(id uint, questionData *models.Question) (*models.Question, error) {
	existingQuestion, err := s.store.Questions.Get(id)
	if err != nil {
		return nil, fmt.Errorf("pregunta no encontrada: %w", err)
	}

	// Update fields
	existingQuestion.Text = questionData.Text
	existingQuestion.Type = questionData.Type
	existingQuestion.Explanation = questionData.Explanation
	existingQuestion.Points = questionData.Points

	if err := s.store.Questions.Update(existingQuestion); err != nil {
		return nil, fmt.Errorf("error al actualizar la pregunta: %w", err)
	}

	return existingQuestion, nil
}

func (s *questionService) UpdateQuestionPatch(questionID uint, data map[string]interface{}) (*models.Question, error) {
	if questionID == 0 {
		return nil, errors.New("question ID cannot be zero")
	}

	var question dto.UpdateQuestionRequest
	if err := utils.MapToStructStrict(data, &question); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.Questions.Patch(questionID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Questions.Get(questionID)
	if err != nil {
		return nil, errors.New("pregunta no encontrada")
	}

	return updated, nil
}

func (s *questionService) DeleteQuestion(id uint) error {
	if err := s.store.Questions.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar la pregunta: %w", err)
	}
	return nil
}

func (s *questionService) GetQuestionsByEvaluation(evaluationID uint) ([]*models.Question, error) {
	// Use the new repository method to filter by evaluation ID at database level
	questions, err := s.store.Questions.GetByEvaluationID(evaluationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las preguntas: %w", err)
	}

	return questions, nil
}

func (s *questionService) GetQuestionWithAnswers(id uint) (*models.Question, error) {
	// Use the new repository method to preload answers
	question, err := s.store.Questions.GetWithAnswers(id)
	if err != nil {
		return nil, fmt.Errorf("pregunta no encontrada: %w", err)
	}

	return question, nil
}

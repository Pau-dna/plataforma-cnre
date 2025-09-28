package services

import (
	"errors"
	"fmt"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type AnswerService interface {
	CreateAnswer(answer *models.Answer) (*models.Answer, error)
	GetAnswer(id uint) (*models.Answer, error)
	UpdateAnswer(id uint, answer *models.Answer) (*models.Answer, error)
	UpdateAnswerPatch(id uint, data map[string]interface{}) (*models.Answer, error)
	DeleteAnswer(id uint) error
	GetAnswersByQuestion(questionID uint) ([]*models.Answer, error)
	ValidateAnswers(questionID uint, selectedAnswerIDs []uint) (bool, int, error)
}

type answerService struct {
	*Service
}

func NewAnswerService(service *Service) AnswerService {
	return &answerService{
		Service: service,
	}
}

func (s *answerService) CreateAnswer(answer *models.Answer) (*models.Answer, error) {
	// Verify question exists
	_, err := s.store.Questions.Get(answer.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("pregunta no encontrada: %w", err)
	}

	if err := s.store.Answers.Create(answer); err != nil {
		return nil, fmt.Errorf("error al crear la respuesta: %w", err)
	}
	return answer, nil
}

func (s *answerService) GetAnswer(id uint) (*models.Answer, error) {
	answer, err := s.store.Answers.Get(id)
	if err != nil {
		return nil, fmt.Errorf("respuesta no encontrada: %w", err)
	}
	return answer, nil
}

func (s *answerService) UpdateAnswer(id uint, answerData *models.Answer) (*models.Answer, error) {
	existingAnswer, err := s.store.Answers.Get(id)
	if err != nil {
		return nil, fmt.Errorf("respuesta no encontrada: %w", err)
	}

	// Update fields
	existingAnswer.Text = answerData.Text
	existingAnswer.IsCorrect = answerData.IsCorrect
	existingAnswer.Order = answerData.Order

	if err := s.store.Answers.Update(existingAnswer); err != nil {
		return nil, fmt.Errorf("error al actualizar la respuesta: %w", err)
	}

	return existingAnswer, nil
}

func (s *answerService) UpdateAnswerPatch(answerID uint, data map[string]interface{}) (*models.Answer, error) {
	if answerID == 0 {
		return nil, errors.New("answer ID cannot be zero")
	}

	var answer dto.UpdateAnswerRequest
	if err := utils.MapToStructStrict(data, &answer); err != nil {
		return nil, errors.New("datos inválidos: " + err.Error())
	}

	if err := s.store.Answers.Patch(answerID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.Answers.Get(answerID)
	if err != nil {
		return nil, errors.New("respuesta no encontrada")
	}

	return updated, nil
}

func (s *answerService) DeleteAnswer(id uint) error {
	if err := s.store.Answers.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar la respuesta: %w", err)
	}
	return nil
}

func (s *answerService) GetAnswersByQuestion(questionID uint) ([]*models.Answer, error) {
	// Use the new repository method to filter by question ID at database level
	answers, err := s.store.Answers.GetByQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las respuestas: %w", err)
	}

	return answers, nil
}

func (s *answerService) ValidateAnswers(questionID uint, selectedAnswerIDs []uint) (bool, int, error) {
	// Get question to determine points
	question, err := s.store.Questions.Get(questionID)
	if err != nil {
		return false, 0, fmt.Errorf("pregunta no encontrada: %w", err)
	}

	// Get all answers for the question
	answers, err := s.GetAnswersByQuestion(questionID)
	if err != nil {
		return false, 0, fmt.Errorf("error al obtener las respuestas: %w", err)
	}

	// Create maps for easy lookup
	answerMap := make(map[uint]*models.Answer)
	correctAnswerIDs := make(map[uint]bool)

	for _, answer := range answers {
		answerMap[answer.ID] = answer
		if answer.IsCorrect {
			correctAnswerIDs[answer.ID] = true
		}
	}

	// Validate selected answers
	selectedMap := make(map[uint]bool)
	for _, id := range selectedAnswerIDs {
		selectedMap[id] = true
	}

	// Check if the selection matches the correct answers exactly
	isCorrect := true

	// Check if all correct answers are selected
	for correctID := range correctAnswerIDs {
		if !selectedMap[correctID] {
			isCorrect = false
			break
		}
	}

	// Check if any incorrect answers are selected
	if isCorrect {
		for _, selectedID := range selectedAnswerIDs {
			if !correctAnswerIDs[selectedID] {
				isCorrect = false
				break
			}
		}
	}

	points := 0
	if isCorrect {
		points = question.Points
	}

	return isCorrect, points, nil
}

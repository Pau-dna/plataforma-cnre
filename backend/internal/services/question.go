package services

import (
	"fmt"

	"github.com/imlargo/go-api-template/internal/models"
)

type QuestionService interface {
	CreateQuestion(question *models.Question) (*models.Question, error)
	GetQuestion(id uint) (*models.Question, error)
	UpdateQuestion(id uint, question *models.Question) (*models.Question, error)
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
		return nil, fmt.Errorf("evaluation not found: %w", err)
	}

	if err := s.store.Questions.Create(question); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}
	return question, nil
}

func (s *questionService) GetQuestion(id uint) (*models.Question, error) {
	question, err := s.store.Questions.Get(id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}
	return question, nil
}

func (s *questionService) UpdateQuestion(id uint, questionData *models.Question) (*models.Question, error) {
	existingQuestion, err := s.store.Questions.Get(id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	// Update fields
	existingQuestion.Text = questionData.Text
	existingQuestion.Type = questionData.Type
	existingQuestion.Explanation = questionData.Explanation
	existingQuestion.Points = questionData.Points

	if err := s.store.Questions.Update(existingQuestion); err != nil {
		return nil, fmt.Errorf("failed to update question: %w", err)
	}

	return existingQuestion, nil
}

func (s *questionService) DeleteQuestion(id uint) error {
	if err := s.store.Questions.Delete(id); err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	return nil
}

func (s *questionService) GetQuestionsByEvaluation(evaluationID uint) ([]*models.Question, error) {
	// This would require a repository method to filter by evaluation ID
	// For now, we'll implement a basic version
	questions, err := s.store.Questions.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	// Filter by evaluation ID
	var evaluationQuestions []*models.Question
	for _, question := range questions {
		if question.EvaluationID == evaluationID {
			evaluationQuestions = append(evaluationQuestions, question)
		}
	}

	return evaluationQuestions, nil
}

func (s *questionService) GetQuestionWithAnswers(id uint) (*models.Question, error) {
	// This would require a repository method to preload answers
	question, err := s.store.Questions.Get(id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}
	
	// For now, return the question - would need to implement preloading in repository
	return question, nil
}
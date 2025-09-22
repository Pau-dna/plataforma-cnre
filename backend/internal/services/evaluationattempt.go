package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/utils"
)

type EvaluationAttemptService interface {
	StartAttempt(userID, evaluationID uint) (*models.EvaluationAttempt, error)
	SubmitAttempt(attemptID uint, answers []models.AttemptAnswer) (*models.EvaluationAttempt, error)
	GetAttempt(id uint) (*models.EvaluationAttempt, error)
	UpdateEvaluationAttemptPatch(id uint, data map[string]interface{}) (*models.EvaluationAttempt, error)
	GetUserAttempts(userID, evaluationID uint) ([]*models.EvaluationAttempt, error)
	CanUserAttempt(userID, evaluationID uint) (bool, string, error)
	ScoreAttempt(attemptID uint) (*models.EvaluationAttempt, error)
}

type evaluationAttemptService struct {
	*Service
	answerService AnswerService
}

func NewEvaluationAttemptService(service *Service, answerService AnswerService) EvaluationAttemptService {
	return &evaluationAttemptService{
		Service:       service,
		answerService: answerService,
	}
}

func (s *evaluationAttemptService) StartAttempt(userID, evaluationID uint) (*models.EvaluationAttempt, error) {
	// Verify user exists
	_, err := s.store.Users.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Verify evaluation exists
	_, err = s.store.Evaluations.Get(evaluationID)
	if err != nil {
		return nil, fmt.Errorf("evaluation not found: %w", err)
	}

	// Check if user can attempt this evaluation
	canAttempt, reason, err := s.CanUserAttempt(userID, evaluationID)
	if err != nil {
		return nil, err
	}
	if !canAttempt {
		return nil, fmt.Errorf("cannot start attempt: %s", reason)
	}

	// Create new attempt
	attempt := &models.EvaluationAttempt{
		UserID:       userID,
		EvaluationID: evaluationID,
		StartedAt:    time.Now(),
		Score:        0,
		TotalPoints:  0,
		Passed:       false,
		Answers:      models.AttemptAnswers{},
	}

	if err := s.store.EvaluationAttempts.Create(attempt); err != nil {
		return nil, fmt.Errorf("failed to create evaluation attempt: %w", err)
	}

	return attempt, nil
}

func (s *evaluationAttemptService) SubmitAttempt(attemptID uint, answers []models.AttemptAnswer) (*models.EvaluationAttempt, error) {
	// Get existing attempt
	attempt, err := s.store.EvaluationAttempts.Get(attemptID)
	if err != nil {
		return nil, fmt.Errorf("attempt not found: %w", err)
	}

	// Check if already submitted
	if !attempt.SubmittedAt.IsZero() {
		return nil, fmt.Errorf("attempt already submitted")
	}

	// Get evaluation to check time limit
	evaluation, err := s.store.Evaluations.Get(attempt.EvaluationID)
	if err != nil {
		return nil, fmt.Errorf("evaluation not found: %w", err)
	}

	// Check time limit if set
	if evaluation.TimeLimit > 0 {
		elapsed := time.Since(attempt.StartedAt)
		if int(elapsed.Minutes()) > evaluation.TimeLimit {
			return nil, fmt.Errorf("time limit exceeded")
		}
	}

	// Calculate time spent
	timeSpent := int(time.Since(attempt.StartedAt).Minutes())

	// Set answers and submission time
	attempt.Answers = models.AttemptAnswers(answers)
	attempt.SubmittedAt = time.Now()
	attempt.TimeSpent = timeSpent

	// Save attempt with answers
	if err := s.store.EvaluationAttempts.Update(attempt); err != nil {
		return nil, fmt.Errorf("failed to update attempt: %w", err)
	}

	// Score the attempt
	scoredAttempt, err := s.ScoreAttempt(attemptID)
	if err != nil {
		return nil, fmt.Errorf("failed to score attempt: %w", err)
	}

	return scoredAttempt, nil
}

func (s *evaluationAttemptService) GetAttempt(id uint) (*models.EvaluationAttempt, error) {
	attempt, err := s.store.EvaluationAttempts.Get(id)
	if err != nil {
		return nil, fmt.Errorf("attempt not found: %w", err)
	}
	return attempt, nil
}

func (s *evaluationAttemptService) UpdateEvaluationAttemptPatch(attemptID uint, data map[string]interface{}) (*models.EvaluationAttempt, error) {
	if attemptID == 0 {
		return nil, errors.New("attempt ID cannot be zero")
	}

	var attempt dto.UpdateEvaluationAttemptRequest
	if err := utils.MapToStructStrict(data, &attempt); err != nil {
		return nil, errors.New("invalid data: " + err.Error())
	}

	if err := s.store.EvaluationAttempts.Patch(attemptID, data); err != nil {
		return nil, err
	}

	updated, err := s.store.EvaluationAttempts.Get(attemptID)
	if err != nil {
		return nil, errors.New("attempt not found")
	}

	return updated, nil
}

func (s *evaluationAttemptService) GetUserAttempts(userID, evaluationID uint) ([]*models.EvaluationAttempt, error) {
	// This would require a repository method to filter by user ID and evaluation ID
	// For now, we'll implement a basic version
	attempts, err := s.store.EvaluationAttempts.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	// Filter by user ID and evaluation ID
	var userAttempts []*models.EvaluationAttempt
	for _, attempt := range attempts {
		if attempt.UserID == userID && attempt.EvaluationID == evaluationID {
			userAttempts = append(userAttempts, attempt)
		}
	}

	return userAttempts, nil
}

func (s *evaluationAttemptService) CanUserAttempt(userID, evaluationID uint) (bool, string, error) {
	// Get evaluation to check max attempts
	evaluation, err := s.store.Evaluations.Get(evaluationID)
	if err != nil {
		return false, "", fmt.Errorf("evaluation not found: %w", err)
	}

	// If no max attempts set, user can always attempt
	if evaluation.MaxAttempts <= 0 {
		return true, "", nil
	}

	// Get user's previous attempts
	attempts, err := s.GetUserAttempts(userID, evaluationID)
	if err != nil {
		return false, "", fmt.Errorf("failed to get user attempts: %w", err)
	}

	// Count completed attempts (submitted)
	completedAttempts := 0
	for _, attempt := range attempts {
		if !attempt.SubmittedAt.IsZero() {
			completedAttempts++
		}
	}

	if completedAttempts >= evaluation.MaxAttempts {
		return false, "maximum attempts reached", nil
	}

	// Check if there's an ongoing attempt
	for _, attempt := range attempts {
		if attempt.SubmittedAt.IsZero() {
			return false, "attempt already in progress", nil
		}
	}

	return true, "", nil
}

func (s *evaluationAttemptService) ScoreAttempt(attemptID uint) (*models.EvaluationAttempt, error) {
	// Get attempt
	attempt, err := s.store.EvaluationAttempts.Get(attemptID)
	if err != nil {
		return nil, fmt.Errorf("attempt not found: %w", err)
	}

	// Get evaluation
	evaluation, err := s.store.Evaluations.Get(attempt.EvaluationID)
	if err != nil {
		return nil, fmt.Errorf("evaluation not found: %w", err)
	}

	totalScore := 0
	totalPoints := 0

	// Score each answer
	for i, answer := range attempt.Answers {
		isCorrect, points, err := s.answerService.ValidateAnswers(
			answer.QuestionID,
			answer.SelectedAnswerIDs,
		)
		if err != nil {
			s.logger.Warnf("Failed to validate answer for question %d: %v", answer.QuestionID, err)
			continue
		}

		// Update the answer in the attempt
		attempt.Answers[i].IsCorrect = isCorrect
		attempt.Answers[i].Points = points

		totalScore += points

		// Get question to add to total points
		question, err := s.store.Questions.Get(answer.QuestionID)
		if err == nil {
			totalPoints += question.Points
		}
	}

	// Update attempt with scores
	attempt.Score = totalScore
	attempt.TotalPoints = totalPoints

	// Check if passed based on passing score
	if totalPoints > 0 {
		percentage := float64(totalScore) / float64(totalPoints) * 100
		attempt.Passed = percentage >= float64(evaluation.PassingScore)
	}

	// Save updated attempt
	if err := s.store.EvaluationAttempts.Update(attempt); err != nil {
		return nil, fmt.Errorf("failed to update attempt scores: %w", err)
	}

	return attempt, nil
}

package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/enums"
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

	// Verify evaluation exists and get configuration
	evaluation, err := s.store.Evaluations.Get(evaluationID)
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

	// Generate dynamic questions and answers for this attempt
	attemptQuestions, totalPoints, err := s.generateAttemptQuestions(evaluation)
	if err != nil {
		return nil, fmt.Errorf("failed to generate attempt questions: %w", err)
	}

	// Create new attempt with generated questions
	attempt := &models.EvaluationAttempt{
		UserID:       userID,
		EvaluationID: evaluationID,
		Questions:    attemptQuestions,
		StartedAt:    time.Now(),
		Score:        0,
		TotalPoints:  totalPoints,
		Passed:       false,
		Answers:      models.AttemptAnswers{},
	}

	if err := s.store.EvaluationAttempts.Create(attempt); err != nil {
		return nil, fmt.Errorf("failed to create evaluation attempt: %w", err)
	}

	return attempt, nil
}

// generateAttemptQuestions generates random questions and answer options for an attempt
func (s *evaluationAttemptService) generateAttemptQuestions(evaluation *models.Evaluation) (models.AttemptQuestions, int, error) {
	// Get all questions for this evaluation
	allQuestions, err := s.store.Questions.GetByEvaluationID(evaluation.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get questions: %w", err)
	}

	if len(allQuestions) < evaluation.QuestionCount {
		return nil, 0, fmt.Errorf("insufficient questions available: need %d, have %d",
			evaluation.QuestionCount, len(allQuestions))
	}

	// Randomly select questions
	selectedQuestions := s.selectRandomQuestions(allQuestions, evaluation.QuestionCount)

	var attemptQuestions models.AttemptQuestions
	totalPoints := 0

	for i, question := range selectedQuestions {
		// Get all answers for this question
		allAnswers := question.Answers

		// Generate answer options for this question
		answerOptions, err := s.generateAnswerOptions(allAnswers, evaluation.AnswerOptionsCount, question.Type)
		if err != nil {
			s.logger.Warnf("Failed to generate answer options for question %d: %v", question.ID, err)
			continue
		}

		attemptQuestion := models.AttemptQuestion{
			ID:            uint(i + 1), // Sequential ID for this attempt
			Text:          question.Text,
			Type:          question.Type,
			Explanation:   question.Explanation,
			Points:        question.Points,
			OriginalID:    question.ID,
			AnswerOptions: answerOptions,
		}

		attemptQuestions = append(attemptQuestions, attemptQuestion)
		totalPoints += question.Points
	}

	if len(attemptQuestions) == 0 {
		return nil, 0, fmt.Errorf("failed to generate any valid questions")
	}

	return attemptQuestions, totalPoints, nil
}

// selectRandomQuestions randomly selects the specified number of questions
func (s *evaluationAttemptService) selectRandomQuestions(questions []*models.Question, count int) []*models.Question {
	if len(questions) <= count {
		return questions
	}

	// Create a copy of the slice to avoid modifying the original
	selected := make([]*models.Question, len(questions))
	copy(selected, questions)

	// Shuffle using Fisher-Yates algorithm
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(selected) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		selected[i], selected[j] = selected[j], selected[i]
	}

	return selected[:count]
}

// generateAnswerOptions generates random answer options ensuring proper distribution
func (s *evaluationAttemptService) generateAnswerOptions(allAnswers []*models.Answer, optionsCount int, questionType enums.QuestionType) ([]models.AttemptAnswerOption, error) {
	if len(allAnswers) < 2 {
		return nil, fmt.Errorf("insufficient answers available: need at least 2, have %d", len(allAnswers))
	}

	// Separate correct and incorrect answers
	var correctAnswers []*models.Answer
	var incorrectAnswers []*models.Answer

	for _, answer := range allAnswers {
		if answer.IsCorrect {
			correctAnswers = append(correctAnswers, answer)
		} else {
			incorrectAnswers = append(incorrectAnswers, answer)
		}
	}

	if len(correctAnswers) == 0 {
		return nil, fmt.Errorf("no correct answers available")
	}

	// Calculate how many correct and incorrect answers to include based on question type
	var maxCorrect int
	var correctToInclude int
	
	switch questionType {
	case enums.QuestionTypeSingle:
		// Single choice questions must have exactly 1 correct answer
		maxCorrect = 1
		correctToInclude = 1
	case enums.QuestionTypeMultiple:
		// Multiple choice questions can have 1 to (optionsCount-1) correct answers
		maxCorrect = optionsCount - 1
		if maxCorrect < 1 {
			maxCorrect = 1
		}
		
		// Random number of correct answers between 1 and maxCorrect
		correctToInclude = 1
		if len(correctAnswers) > 1 && maxCorrect > 1 {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			correctToInclude = r.Intn(maxCorrect) + 1
		}
	default:
		// Fallback to original logic for unknown question types
		maxCorrect = optionsCount / 2
		if maxCorrect == 0 {
			maxCorrect = 1 // Always at least 1 correct answer
		}
		
		correctToInclude = 1
		if len(correctAnswers) > 1 && maxCorrect > 1 {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			correctToInclude = r.Intn(maxCorrect) + 1
		}
	}

	// Don't exceed available correct answers
	if correctToInclude > len(correctAnswers) {
		correctToInclude = len(correctAnswers)
	}

	incorrectToInclude := optionsCount - correctToInclude
	if incorrectToInclude > len(incorrectAnswers) {
		incorrectToInclude = len(incorrectAnswers)
	}

	// If we don't have enough incorrect answers, adjust
	if incorrectToInclude < 0 {
		incorrectToInclude = 0
	}

	// Select random correct and incorrect answers
	selectedCorrect := s.selectRandomAnswers(correctAnswers, correctToInclude)
	selectedIncorrect := s.selectRandomAnswers(incorrectAnswers, incorrectToInclude)

	// Combine and shuffle the options
	var options []models.AttemptAnswerOption

	// Add correct options
	for i, answer := range selectedCorrect {
		options = append(options, models.AttemptAnswerOption{
			ID:        uint(i + 1),
			Text:      answer.Text,
			IsCorrect: true,
		})
	}

	// Add incorrect options
	for i, answer := range selectedIncorrect {
		options = append(options, models.AttemptAnswerOption{
			ID:        uint(len(selectedCorrect) + i + 1),
			Text:      answer.Text,
			IsCorrect: false,
		})
	}

	// Shuffle the options so correct answers aren't always first
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(options) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	// Reassign sequential IDs after shuffling
	for i := range options {
		options[i].ID = uint(i + 1)
	}

	return options, nil
}

// selectRandomAnswers randomly selects the specified number of answers
func (s *evaluationAttemptService) selectRandomAnswers(answers []*models.Answer, count int) []*models.Answer {
	if len(answers) <= count {
		return answers
	}

	// Create a copy to avoid modifying the original
	selected := make([]*models.Answer, len(answers))
	copy(selected, answers)

	// Shuffle
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(selected) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		selected[i], selected[j] = selected[j], selected[i]
	}

	return selected[:count]
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

	// Get evaluation to check time limit (cache this if needed)
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

	// Perform scoring and saving in a single optimized operation
	s.scoreAttemptInline(attempt, evaluation)

	// Save attempt with all updates (answers + scores) in single transaction
	if err := s.store.EvaluationAttempts.Update(attempt); err != nil {
		return nil, fmt.Errorf("failed to update attempt: %w", err)
	}

	return attempt, nil
}

// scoreAttemptInline performs scoring directly on the attempt object without database round trips
func (s *evaluationAttemptService) scoreAttemptInline(attempt *models.EvaluationAttempt, evaluation *models.Evaluation) {
	totalScore := 0
	totalPoints := attempt.TotalPoints

	// Create a map for faster question lookups
	questionMap := make(map[uint]*models.AttemptQuestion)
	for i := range attempt.Questions {
		questionMap[attempt.Questions[i].ID] = &attempt.Questions[i]
	}

	// Score each answer using the optimized map lookup
	for i, answer := range attempt.Answers {
		// Find the corresponding attempt question using map (O(1) vs O(n))
		attemptQuestion, exists := questionMap[answer.AttemptQuestionID]
		if !exists {
			s.logger.Warnf("Attempt question %d not found for answer", answer.AttemptQuestionID)
			continue
		}

		// Validate the selected options
		isCorrect, points := s.validateAttemptAnswer(attemptQuestion, answer.SelectedOptionIDs)

		// Update the answer in the attempt
		attempt.Answers[i].IsCorrect = isCorrect
		attempt.Answers[i].Points = points

		totalScore += points
	}

	// Update attempt with scores
	attempt.Score = totalScore

	// Check if passed based on passing score
	if totalPoints > 0 {
		percentage := float64(totalScore) / float64(totalPoints) * 100
		attempt.Passed = percentage >= float64(evaluation.PassingScore)
	}
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
	// Use the optimized repository method for efficient database-level filtering
	attempts, err := s.store.EvaluationAttempts.GetByUserAndEvaluation(userID, evaluationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	return attempts, nil
}

func (s *evaluationAttemptService) CanUserAttempt(userID, evaluationID uint) (bool, string, error) {
	// Get evaluation to check max attempts
	evaluation, err := s.store.Evaluations.Get(evaluationID)
	if err != nil {
		return false, "", fmt.Errorf("evaluation not found: %w", err)
	}

	// If no max attempts set, user can always attempt
	if evaluation.MaxAttempts <= 0 {
		// Still check for in-progress attempts
		_, err := s.store.EvaluationAttempts.GetInProgressAttempt(userID, evaluationID)
		if err == nil {
			return false, "attempt already in progress", nil
		}
		return true, "", nil
	}

	// Use optimized database query to count completed attempts
	completedAttempts, err := s.store.EvaluationAttempts.CountCompletedAttempts(userID, evaluationID)
	if err != nil {
		return false, "", fmt.Errorf("failed to count user attempts: %w", err)
	}

	if int(completedAttempts) >= evaluation.MaxAttempts {
		return false, "maximum attempts reached", nil
	}

	// Check if there's an ongoing attempt with optimized query
	_, err = s.store.EvaluationAttempts.GetInProgressAttempt(userID, evaluationID)
	if err == nil {
		return false, "attempt already in progress", nil
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
	totalPoints := attempt.TotalPoints // Already calculated during attempt creation

	// Create a map for faster question lookups instead of nested loops
	questionMap := make(map[uint]*models.AttemptQuestion)
	for i := range attempt.Questions {
		questionMap[attempt.Questions[i].ID] = &attempt.Questions[i]
	}

	// Score each answer using the optimized map lookup
	for i, answer := range attempt.Answers {
		// Find the corresponding attempt question using map (O(1) vs O(n))
		attemptQuestion, exists := questionMap[answer.AttemptQuestionID]
		if !exists {
			s.logger.Warnf("Attempt question %d not found for answer", answer.AttemptQuestionID)
			continue
		}

		// Validate the selected options
		isCorrect, points := s.validateAttemptAnswer(attemptQuestion, answer.SelectedOptionIDs)

		// Update the answer in the attempt
		attempt.Answers[i].IsCorrect = isCorrect
		attempt.Answers[i].Points = points

		totalScore += points
	}

	// Update attempt with scores
	attempt.Score = totalScore

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

// validateAttemptAnswer validates user's selected options against the attempt question
func (s *evaluationAttemptService) validateAttemptAnswer(question *models.AttemptQuestion, selectedOptionIDs []uint) (bool, int) {
	if len(selectedOptionIDs) == 0 {
		return false, 0
	}

	// Get correct option IDs
	var correctOptionIDs []uint
	for _, option := range question.AnswerOptions {
		if option.IsCorrect {
			correctOptionIDs = append(correctOptionIDs, option.ID)
		}
	}

	// Check if all selected options are correct and all correct options are selected
	selectedMap := make(map[uint]bool)
	for _, id := range selectedOptionIDs {
		selectedMap[id] = true
	}

	correctMap := make(map[uint]bool)
	for _, id := range correctOptionIDs {
		correctMap[id] = true
	}

	// For single choice: only one option should be selected and it should be correct
	if question.Type == "single_choice" {
		if len(selectedOptionIDs) != 1 {
			return false, 0
		}
		if correctMap[selectedOptionIDs[0]] {
			return true, question.Points
		}
		return false, 0
	}

	// For multiple choice: all selected must be correct, and all correct must be selected
	if question.Type == "multiple_choice" {
		// Check if all selected are correct
		for _, id := range selectedOptionIDs {
			if !correctMap[id] {
				return false, 0 // Selected an incorrect option
			}
		}

		// Check if all correct options are selected
		for _, id := range correctOptionIDs {
			if !selectedMap[id] {
				return false, 0 // Missed a correct option
			}
		}

		return true, question.Points
	}

	return false, 0
}

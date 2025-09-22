package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/internal/services"
)

type EvaluationAttemptHandler struct {
	*Handler
	evaluationAttemptService services.EvaluationAttemptService
}

func NewEvaluationAttemptHandler(handler *Handler, evaluationAttemptService services.EvaluationAttemptService) *EvaluationAttemptHandler {
	return &EvaluationAttemptHandler{
		Handler:                  handler,
		evaluationAttemptService: evaluationAttemptService,
	}
}

// @Summary Start evaluation attempt
// @Description Start a new evaluation attempt for a user
// @Tags evaluation-attempts
// @Accept json
// @Produce json
// @Success 201 {object} models.EvaluationAttempt
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluation-attempts/start [post]
func (h *EvaluationAttemptHandler) StartAttempt(c *gin.Context) {
	var attemptData struct {
		UserID       uint `json:"user_id" binding:"required"`
		EvaluationID uint `json:"evaluation_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&attemptData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt, err := h.evaluationAttemptService.StartAttempt(attemptData.UserID, attemptData.EvaluationID)
	if err != nil {
		h.logger.Errorf("Failed to start attempt: %v", err)
		if err.Error() == "cannot start attempt: maximum attempts reached" ||
			err.Error() == "cannot start attempt: attempt already in progress" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start attempt"})
		return
	}

	c.JSON(http.StatusCreated, attempt)
}

// @Summary Submit evaluation attempt
// @Description Submit answers for an evaluation attempt
// @Tags evaluation-attempts
// @Accept json
// @Produce json
// @Param id path int true "Attempt ID"
// @Success 200 {object} models.EvaluationAttempt
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluation-attempts/{id}/submit [post]
func (h *EvaluationAttemptHandler) SubmitAttempt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	var submissionData struct {
		Answers []models.AttemptAnswer `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&submissionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt, err := h.evaluationAttemptService.SubmitAttempt(uint(id), submissionData.Answers)
	if err != nil {
		h.logger.Errorf("Failed to submit attempt: %v", err)
		if err.Error() == "attempt already submitted" || err.Error() == "time limit exceeded" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit attempt"})
		return
	}

	c.JSON(http.StatusOK, attempt)
}

// @Summary Get evaluation attempt
// @Description Get an evaluation attempt by its ID
// @Tags evaluation-attempts
// @Produce json
// @Param id path int true "Attempt ID"
// @Success 200 {object} models.EvaluationAttempt
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/evaluation-attempts/{id} [get]
func (h *EvaluationAttemptHandler) GetAttempt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	attempt, err := h.evaluationAttemptService.GetAttempt(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get attempt: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Attempt not found"})
		return
	}

	c.JSON(http.StatusOK, attempt)
}

// @Summary Get user evaluation attempts
// @Description Get all attempts for a user and evaluation
// @Tags evaluation-attempts
// @Produce json
// @Param userId path int true "User ID"
// @Param evaluationId path int true "Evaluation ID"
// @Success 200 {array} models.EvaluationAttempt
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/evaluations/{evaluationId}/attempts [get]
func (h *EvaluationAttemptHandler) GetUserAttempts(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	evaluationIDStr := c.Param("evaluationId")
	evaluationID, err := strconv.ParseUint(evaluationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation ID"})
		return
	}

	attempts, err := h.evaluationAttemptService.GetUserAttempts(uint(userID), uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Failed to get user attempts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attempts"})
		return
	}

	c.JSON(http.StatusOK, attempts)
}

// @Summary Check if user can attempt evaluation
// @Description Check if a user can attempt an evaluation
// @Tags evaluation-attempts
// @Produce json
// @Param userId path int true "User ID"
// @Param evaluationId path int true "Evaluation ID"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/evaluations/{evaluationId}/can-attempt [get]
func (h *EvaluationAttemptHandler) CanUserAttempt(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	evaluationIDStr := c.Param("evaluationId")
	evaluationID, err := strconv.ParseUint(evaluationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation ID"})
		return
	}

	canAttempt, reason, err := h.evaluationAttemptService.CanUserAttempt(uint(userID), uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Failed to check if user can attempt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check attempt eligibility"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"can_attempt": canAttempt,
		"reason":      reason,
	})
}

// @Summary Score evaluation attempt
// @Description Score an evaluation attempt (admin only)
// @Tags evaluation-attempts
// @Param id path int true "Attempt ID"
// @Success 200 {object} models.EvaluationAttempt
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluation-attempts/{id}/score [post]
func (h *EvaluationAttemptHandler) ScoreAttempt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	attempt, err := h.evaluationAttemptService.ScoreAttempt(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to score attempt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to score attempt"})
		return
	}

	c.JSON(http.StatusOK, attempt)
}

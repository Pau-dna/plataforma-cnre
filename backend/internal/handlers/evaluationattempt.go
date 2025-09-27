package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"

	"github.com/imlargo/go-api-template/internal/responses"
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
		responses.ErrorBindJson(c, err)
		return
	}

	attempt, err := h.evaluationAttemptService.StartAttempt(attemptData.UserID, attemptData.EvaluationID)
	if err != nil {
		h.logger.Errorf("Failed to start attempt: %v", err)
		if err.Error() == "cannot start attempt: maximum attempts reached" ||
			err.Error() == "cannot start attempt: attempt already in progress" {
			responses.ErrorConflict(c, err.Error())
			return
		}
		responses.ErrorInternalServerWithMessage(c, "Failed to start attempt:"+err.Error())
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
		responses.ErrorBadRequest(c, "Invalid attempt ID")
		return
	}

	var submissionData struct {
		Answers []models.AttemptAnswer `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&submissionData); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	attempt, err := h.evaluationAttemptService.SubmitAttempt(uint(id), submissionData.Answers)
	if err != nil {
		h.logger.Errorf("Failed to submit attempt: %v", err)
		if err.Error() == "attempt already submitted" || err.Error() == "time limit exceeded" {
			responses.ErrorConflict(c, err.Error())
			return
		}
		responses.ErrorInternalServerWithMessage(c, "Failed to submit attempt")
		return
	}

	responses.Ok(c, attempt)
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
		responses.ErrorBadRequest(c, "Invalid attempt ID")
		return
	}

	attempt, err := h.evaluationAttemptService.GetAttempt(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get attempt: %v", err)
		responses.ErrorNotFound(c, "Attempt")
		return
	}

	responses.Ok(c, attempt)
}

// @Summary		Update evaluation attempt
// @Router			/api/v1/evaluation-attempts/{id} [patch]
// @Description	Update an evaluation attempt by ID
// @Tags		evaluation-attempts
// @Param id path int true "Attempt ID"
// @Accept		json
// @Param payload body dto.UpdateEvaluationAttemptRequest true "Attempt data"
// @Produce		json
// @Success		200	{object}	models.EvaluationAttempt	"Attempt updated successfully"
// @Failure		400	{object}	responses.ErrorResponse	"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse	"Attempt not found"
// @Failure		500	{object}	responses.ErrorResponse	"Internal Server Error"
// @Security     BearerAuth
func (h *EvaluationAttemptHandler) UpdateEvaluationAttemptPatch(c *gin.Context) {
	attemptID := c.Param("id")
	if attemptID == "" {
		responses.ErrorBadRequest(c, "Attempt ID is required")
		return
	}

	attemptIDInt, err := strconv.Atoi(attemptID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Attempt ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	attempt, err := h.evaluationAttemptService.UpdateEvaluationAttemptPatch(uint(attemptIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, attempt)
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
		responses.ErrorBadRequest(c, "Invalid user ID")
		return
	}

	evaluationIDStr := c.Param("evaluationId")
	evaluationID, err := strconv.ParseUint(evaluationIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid evaluation ID")
		return
	}

	attempts, err := h.evaluationAttemptService.GetUserAttempts(uint(userID), uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Failed to get user attempts: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to get attempts")
		return
	}

	responses.Ok(c, attempts)
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
		responses.ErrorBadRequest(c, "Invalid user ID")
		return
	}

	evaluationIDStr := c.Param("evaluationId")
	evaluationID, err := strconv.ParseUint(evaluationIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid evaluation ID")
		return
	}

	canAttempt, reason, err := h.evaluationAttemptService.CanUserAttempt(uint(userID), uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Failed to check if user can attempt: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to check attempt eligibility")
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
		responses.ErrorBadRequest(c, "Invalid attempt ID")
		return
	}

	attempt, err := h.evaluationAttemptService.ScoreAttempt(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to score attempt: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to score attempt")
		return
	}

	responses.Ok(c, attempt)
}

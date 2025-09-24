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

type EvaluationHandler struct {
	*Handler
	evaluationService services.EvaluationService
}

func NewEvaluationHandler(handler *Handler, evaluationService services.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{
		Handler:           handler,
		evaluationService: evaluationService,
	}
}

// @Summary Create evaluation
// @Description Create new evaluation for a module
// @Tags evaluations
// @Accept json
// @Produce json
// @Param evaluation body models.Evaluation true "Evaluation data"
// @Success 201 {object} models.Evaluation
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluations [post]
func (h *EvaluationHandler) CreateEvaluation(c *gin.Context) {
	var evaluation models.Evaluation
	if err := c.ShouldBindJSON(&evaluation); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	createdEvaluation, err := h.evaluationService.CreateEvaluation(&evaluation)
	if err != nil {
		h.logger.Errorf("Failed to create evaluation: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to create evaluation")
		return
	}

	c.JSON(http.StatusCreated, createdEvaluation)
}

// @Summary Get evaluation
// @Description Get evaluation by ID
// @Tags evaluations
// @Produce json
// @Param id path int true "Evaluation ID"
// @Success 200 {object} models.Evaluation
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/evaluations/{id} [get]
func (h *EvaluationHandler) GetEvaluation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid evaluation ID")
		return
	}

	evaluation, err := h.evaluationService.GetEvaluation(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get evaluation: %v", err)
		responses.ErrorNotFound(c, "Evaluation")
		return
	}

	responses.Ok(c, evaluation)
}

// @Summary Update evaluation
// @Description Update evaluation by ID
// @Tags evaluations
// @Accept json
// @Produce json
// @Param id path int true "Evaluation ID"
// @Param evaluation body models.Evaluation true "Evaluation data"
// @Success 200 {object} models.Evaluation
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluations/{id} [put]
func (h *EvaluationHandler) UpdateEvaluation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid evaluation ID")
		return
	}

	var evaluation models.Evaluation
	if err := c.ShouldBindJSON(&evaluation); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	updatedEvaluation, err := h.evaluationService.UpdateEvaluation(uint(id), &evaluation)
	if err != nil {
		h.logger.Errorf("Failed to update evaluation: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to update evaluation")
		return
	}

	responses.Ok(c, updatedEvaluation)
}

// @SummaryUpdate evaluation
// @Router/api/v1/evaluations/{id} [patch]
// @DescriptionUpdate a evaluation by ID
// @Tagsevaluations
// @Param id path int true "Evaluation ID"
// @Acceptjson
// @Param payload body dto.UpdateEvaluationRequest true "Evaluation data"
// @Producejson
// @Success200{object}models.Evaluation"Evaluation updated successfully"
// @Failure400{object}responses.ErrorResponse"Bad Request"
// @Failure404{object}responses.ErrorResponse"Evaluation not found"
// @Failure500{object}responses.ErrorResponse"Internal Server Error"
// @Security     BearerAuth
func (h *EvaluationHandler) UpdateEvaluationPatch(c *gin.Context) {
	evaluationID := c.Param("id")
	if evaluationID == "" {
		responses.ErrorBadRequest(c, "Evaluation ID is required")
		return
	}

	evaluationIDInt, err := strconv.Atoi(evaluationID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Evaluation ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	evaluation, err := h.evaluationService.UpdateEvaluationPatch(uint(evaluationIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, evaluation)
}

// @Summary Delete evaluation
// @Description Delete evaluation by ID
// @Tags evaluations
// @Param id path int true "Evaluation ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluations/{id} [delete]
func (h *EvaluationHandler) DeleteEvaluation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid evaluation ID")
		return
	}

	err = h.evaluationService.DeleteEvaluation(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete evaluation: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to delete evaluation")
		return
	}

	responses.Ok(c, "ok")
}

// @Summary Get module evaluations
// @Description Get all evaluations for a specific module
// @Tags evaluations
// @Produce json
// @Param moduleId path int true "Module ID"
// @Success 200 {array} models.Evaluation
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/modules/{moduleId}/evaluations [get]
func (h *EvaluationHandler) GetEvaluationsByModule(c *gin.Context) {
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid module ID")
		return
	}

	evaluations, err := h.evaluationService.GetEvaluationsByModule(uint(moduleID))
	if err != nil {
		h.logger.Errorf("Failed to get evaluations by module: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to get evaluations")
		return
	}

	responses.Ok(c, evaluations)
}

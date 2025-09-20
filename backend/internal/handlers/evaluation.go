package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdEvaluation, err := h.evaluationService.CreateEvaluation(&evaluation)
	if err != nil {
		h.logger.Errorf("Failed to create evaluation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create evaluation"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation ID"})
		return
	}

	evaluation, err := h.evaluationService.GetEvaluation(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get evaluation: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Evaluation not found"})
		return
	}

	c.JSON(http.StatusOK, evaluation)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation ID"})
		return
	}

	var evaluation models.Evaluation
	if err := c.ShouldBindJSON(&evaluation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvaluation, err := h.evaluationService.UpdateEvaluation(uint(id), &evaluation)
	if err != nil {
		h.logger.Errorf("Failed to update evaluation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update evaluation"})
		return
	}

	c.JSON(http.StatusOK, updatedEvaluation)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation ID"})
		return
	}

	err = h.evaluationService.DeleteEvaluation(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete evaluation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete evaluation"})
		return
	}

	c.Status(http.StatusNoContent)
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
	moduleIDStr := c.Param("moduleId")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	evaluations, err := h.evaluationService.GetEvaluationsByModule(uint(moduleID))
	if err != nil {
		h.logger.Errorf("Failed to get evaluations by module: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get evaluations"})
		return
	}

	c.JSON(http.StatusOK, evaluations)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/internal/responses"
	"github.com/imlargo/go-api-template/internal/services"
)

type AnswerHandler struct {
	*Handler
	answerService services.AnswerService
}

func NewAnswerHandler(handler *Handler, answerService services.AnswerService) *AnswerHandler {
	return &AnswerHandler{
		Handler:       handler,
		answerService: answerService,
	}
}

// @Summary Create answer
// @Description Create new answer for a question
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body dto.CreateAnswerRequest true "Answer data"
// @Success 201 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/answers [post]
func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	var answerReq dto.CreateAnswerRequest
	if err := c.ShouldBindJSON(&answerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := &models.Answer{
		Text:       answerReq.Text,
		IsCorrect:  answerReq.IsCorrect,
		Order:      answerReq.Order,
		QuestionID: answerReq.QuestionID,
	}

	createdAnswer, err := h.answerService.CreateAnswer(answer)
	if err != nil {
		h.logger.Errorf("Failed to create answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, createdAnswer)
}

// @Summary Get answer
// @Description Get answer by ID
// @Tags answers
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/answers/{id} [get]
func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	answer, err := h.answerService.GetAnswer(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get answer: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, answer)
}

// @Summary Update answer
// @Description Update answer by ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Param answer body dto.UpdateAnswerRequest true "Answer data"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/answers/{id} [put]
func (h *AnswerHandler) UpdateAnswer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	var answerReq dto.UpdateAnswerRequest
	if err := c.ShouldBindJSON(&answerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := &models.Answer{
		Text:      answerReq.Text,
		IsCorrect: answerReq.IsCorrect,
		Order:     answerReq.Order,
	}

	updatedAnswer, err := h.answerService.UpdateAnswer(uint(id), answer)
	if err != nil {
		h.logger.Errorf("Failed to update answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer"})
		return
	}

	c.JSON(http.StatusOK, updatedAnswer)
}

// @Summary Update answer
// @Router /api/v1/answers/{id} [patch]
// @Description Update an answer by ID
// @Tags answers
// @Param id path int true "Answer ID"
// @Accept json
// @Param payload body dto.UpdateAnswerRequest true "Answer data"
// @Produce json
// @Success 200 {object} models.Answer "Answer updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 404 {object} responses.ErrorResponse "Answer not found"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Security BearerAuth
func (h *AnswerHandler) UpdateAnswerPatch(c *gin.Context) {
	answerID := c.Param("id")
	if answerID == "" {
		responses.ErrorBadRequest(c, "Answer ID is required")
		return
	}

	answerIDInt, err := strconv.Atoi(answerID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Answer ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	answer, err := h.answerService.UpdateAnswerPatch(uint(answerIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, answer)
}

// @Summary Delete answer
// @Description Delete answer by ID
// @Tags answers
// @Param id path int true "Answer ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/answers/{id} [delete]
func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	err = h.answerService.DeleteAnswer(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer"})
		return
	}

	responses.Ok(c, "ok")
}

// @Summary Get question answers
// @Description Get all answers for a specific question
// @Tags answers
// @Produce json
// @Param questionId path int true "Question ID"
// @Success 200 {array} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/questions/{questionId}/answers [get]
func (h *AnswerHandler) GetAnswersByQuestion(c *gin.Context) {
	questionIDStr := c.Param("id")
	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	answers, err := h.answerService.GetAnswersByQuestion(uint(questionID))
	if err != nil {
		h.logger.Errorf("Failed to get answers by question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get answers"})
		return
	}

	c.JSON(http.StatusOK, answers)
}
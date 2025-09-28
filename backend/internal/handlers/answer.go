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
		responses.ErrorBindJson(c, err)
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
		responses.ErrorInternalServerWithMessage(c, "Error al crear la respuesta")
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
		responses.ErrorBadRequest(c, "ID de respuesta inválido")
		return
	}

	answer, err := h.answerService.GetAnswer(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get answer: %v", err)
		responses.ErrorNotFound(c, "Respuesta")
		return
	}

	responses.Ok(c, answer)
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
		responses.ErrorBadRequest(c, "ID de respuesta inválido")
		return
	}

	var answerReq dto.UpdateAnswerRequest
	if err := c.ShouldBindJSON(&answerReq); err != nil {
		responses.ErrorBindJson(c, err)
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
		responses.ErrorInternalServerWithMessage(c, "Error al actualizar la respuesta")
		return
	}

	responses.Ok(c, updatedAnswer)
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
		responses.ErrorBadRequest(c, "El ID de respuesta es requerido")
		return
	}

	answerIDInt, err := strconv.Atoi(answerID)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de respuesta inválido: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Datos de la petición inválidos: "+err.Error())
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
		responses.ErrorBadRequest(c, "ID de respuesta inválido")
		return
	}

	err = h.answerService.DeleteAnswer(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete answer: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al eliminar la respuesta")
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
		responses.ErrorBadRequest(c, "ID de pregunta inválido")
		return
	}

	answers, err := h.answerService.GetAnswersByQuestion(uint(questionID))
	if err != nil {
		h.logger.Errorf("Failed to get answers by question: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al obtener las respuestas")
		return
	}

	responses.Ok(c, answers)
}

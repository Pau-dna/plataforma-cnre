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

type QuestionHandler struct {
	*Handler
	questionService services.QuestionService
}

func NewQuestionHandler(handler *Handler, questionService services.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		Handler:         handler,
		questionService: questionService,
	}
}

// @Summary Create question
// @Description Create new question for an evaluation
// @Tags questions
// @Accept json
// @Produce json
// @Param question body dto.CreateQuestionRequest true "Question data"
// @Success 201 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/questions [post]
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var questionReq dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&questionReq); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	question := &models.Question{
		Text:         questionReq.Text,
		Type:         questionReq.Type,
		Explanation:  questionReq.Explanation,
		Points:       questionReq.Points,
		EvaluationID: questionReq.EvaluationID,
	}

	createdQuestion, err := h.questionService.CreateQuestion(question)
	if err != nil {
		h.logger.Errorf("Error al crear la pregunta: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al crear la pregunta")
		return
	}

	c.JSON(http.StatusCreated, createdQuestion)
}

// @Summary Get question
// @Description Get question by ID
// @Tags questions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/questions/{id} [get]
func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de pregunta inválido")
		return
	}

	question, err := h.questionService.GetQuestion(uint(id))
	if err != nil {
		h.logger.Errorf("Error al obtener la pregunta: %v", err)
		responses.ErrorNotFound(c, "Pregunta")
		return
	}

	responses.Ok(c, question)
}

// @Summary Update question
// @Description Update question by ID
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param question body dto.UpdateQuestionRequest true "Question data"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/questions/{id} [put]
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de pregunta inválido")
		return
	}

	var questionReq dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&questionReq); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	question := &models.Question{
		Text:        questionReq.Text,
		Type:        questionReq.Type,
		Explanation: questionReq.Explanation,
		Points:      questionReq.Points,
	}

	updatedQuestion, err := h.questionService.UpdateQuestion(uint(id), question)
	if err != nil {
		h.logger.Errorf("Error al actualizar la pregunta: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al actualizar la pregunta")
		return
	}

	responses.Ok(c, updatedQuestion)
}

// @Summary Update question
// @Router /api/v1/questions/{id} [patch]
// @Description Update a question by ID
// @Tags questions
// @Param id path int true "Question ID"
// @Accept json
// @Param payload body dto.UpdateQuestionRequest true "Question data"
// @Produce json
// @Success 200 {object} models.Question "Question updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 404 {object} responses.ErrorResponse "Question not found"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Security BearerAuth
func (h *QuestionHandler) UpdateQuestionPatch(c *gin.Context) {
	questionID := c.Param("id")
	if questionID == "" {
		responses.ErrorBadRequest(c, "Question ID is required")
		return
	}

	questionIDInt, err := strconv.Atoi(questionID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Question ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	question, err := h.questionService.UpdateQuestionPatch(uint(questionIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, question)
}

// @Summary Delete question
// @Description Delete question by ID
// @Tags questions
// @Param id path int true "Question ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/questions/{id} [delete]
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de pregunta inválido")
		return
	}

	err = h.questionService.DeleteQuestion(uint(id))
	if err != nil {
		h.logger.Errorf("Error al eliminar la pregunta: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al eliminar la pregunta")
		return
	}

	responses.Ok(c, "ok")
}

// @Summary Get evaluation questions
// @Description Get all questions for a specific evaluation
// @Tags questions
// @Produce json
// @Param evaluationId path int true "Evaluation ID"
// @Success 200 {array} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/evaluations/{evaluationId}/questions [get]
func (h *QuestionHandler) GetQuestionsByEvaluation(c *gin.Context) {
	evaluationIDStr := c.Param("id")
	evaluationID, err := strconv.ParseUint(evaluationIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de evaluación inválido")
		return
	}

	questions, err := h.questionService.GetQuestionsByEvaluation(uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Error al obtener la preguntas by evaluation: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al obtener la preguntas")
		return
	}

	responses.Ok(c, questions)
}

// @Summary Get question with answers
// @Description Get question by ID with all its answers
// @Tags questions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/questions/{id}/answers [get]
func (h *QuestionHandler) GetQuestionWithAnswers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de pregunta inválido")
		return
	}

	question, err := h.questionService.GetQuestionWithAnswers(uint(id))
	if err != nil {
		h.logger.Errorf("Error al obtener la pregunta with answers: %v", err)
		responses.ErrorNotFound(c, "Pregunta")
		return
	}

	responses.Ok(c, question)
}

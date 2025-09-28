package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/responses"
	"github.com/imlargo/go-api-template/internal/services"
)

type UserProgressHandler struct {
	*Handler
	userProgressService services.UserProgressService
}

func NewUserProgressHandler(handler *Handler, userProgressService services.UserProgressService) *UserProgressHandler {
	return &UserProgressHandler{
		Handler:             handler,
		userProgressService: userProgressService,
	}
}

// @Summary Mark content as completed
// @Description Mark a specific content item as completed for a user
// @Tags user-progress
// @Accept json
// @Produce json
// @Param data body object true "Content completion data"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-progress/complete [post]
func (h *UserProgressHandler) MarkContentComplete(c *gin.Context) {
	var req struct {
		UserID    uint `json:"user_id" binding:"required"`
		CourseID  uint `json:"course_id" binding:"required"`
		ModuleID  uint `json:"module_id" binding:"required"`
		ContentID uint `json:"content_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	progress, err := h.userProgressService.MarkContentComplete(req.UserID, req.CourseID, req.ModuleID, req.ContentID)
	if err != nil {
		h.logger.Errorf("Error al marcar contenido como completado: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo marcar el contenido como completado")
		return
	}

	c.JSON(http.StatusCreated, progress)
}

// @Summary Mark content as incomplete
// @Description Mark a specific content item as incomplete for a user
// @Tags user-progress
// @Accept json
// @Produce json
// @Param data body object true "Content incompletion data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-progress/incomplete [post]
func (h *UserProgressHandler) MarkContentIncomplete(c *gin.Context) {
	var req struct {
		UserID    uint `json:"user_id" binding:"required"`
		CourseID  uint `json:"course_id" binding:"required"`
		ModuleID  uint `json:"module_id" binding:"required"`
		ContentID uint `json:"content_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	err := h.userProgressService.MarkContentIncomplete(req.UserID, req.CourseID, req.ModuleID, req.ContentID)
	if err != nil {
		h.logger.Errorf("Error al marcar contenido como incompleto: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo marcar el contenido como incompleto")
		return
	}

	responses.Ok(c, gin.H{"message": "Contenido marcado como incompleto"})
}

// @Summary Get user course progress
// @Description Get all progress records for a user in a specific course
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param courseId path int true "Course ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{courseId}/progress [get]
func (h *UserProgressHandler) GetUserCourseProgress(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	progress, err := h.userProgressService.GetUserProgress(uint(userID), uint(courseID))
	if err != nil {
		h.logger.Errorf("Error al obtener progreso del curso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo obtener el progreso del curso")
		return
	}

	responses.Ok(c, progress)
}

// @Summary Get user module progress
// @Description Get all progress records for a user in a specific module
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param moduleId path int true "Module ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/modules/{moduleId}/progress [get]
func (h *UserProgressHandler) GetUserModuleProgress(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de módulo inválido")
		return
	}

	progress, err := h.userProgressService.GetUserModuleProgress(uint(userID), uint(moduleID))
	if err != nil {
		h.logger.Errorf("Error al obtener progreso del módulo: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo obtener el progreso del módulo")
		return
	}

	responses.Ok(c, progress)
}

// @Summary Calculate course progress percentage
// @Description Calculate the progress percentage for a user in a specific course
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} map[string]interface{} "Progress percentage"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{courseId}/progress-percentage [get]
func (h *UserProgressHandler) CalculateCourseProgress(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	progressPercentage, err := h.userProgressService.CalculateCourseProgress(uint(userID), uint(courseID))
	if err != nil {
		h.logger.Errorf("Error al calcular progreso del curso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo calcular el progreso del curso")
		return
	}

	responses.Ok(c, gin.H{"progress_percentage": progressPercentage})
}

// @Summary Calculate module progress percentage
// @Description Calculate the progress percentage for a user in a specific module
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param moduleId path int true "Module ID"
// @Success 200 {object} map[string]interface{} "Progress percentage"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/modules/{moduleId}/progress-percentage [get]
func (h *UserProgressHandler) CalculateModuleProgress(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de módulo inválido")
		return
	}

	progressPercentage, err := h.userProgressService.CalculateModuleProgress(uint(userID), uint(moduleID))
	if err != nil {
		h.logger.Errorf("Error al calcular progreso del módulo: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo calcular el progreso del módulo")
		return
	}

	responses.Ok(c, gin.H{"progress_percentage": progressPercentage})
}

// @Summary Get user progress for specific content
// @Description Get progress record for a user and specific content
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param contentId path int true "Content ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/content/{contentId}/progress [get]
func (h *UserProgressHandler) GetUserContentProgress(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de contenido inválido")
		return
	}

	progress, err := h.userProgressService.GetUserProgressForContent(uint(userID), uint(contentID))
	if err != nil {
		h.logger.Errorf("Error al obtener progreso del contenido: %v", err)
		responses.ErrorNotFound(c, "Progreso de contenido")
		return
	}

	responses.Ok(c, progress)
}

// @Summary Check if user has passed evaluation
// @Description Check if a user has passed a specific evaluation
// @Tags user-progress
// @Produce json
// @Param userId path int true "User ID"
// @Param evaluationId path int true "Evaluation ID"
// @Success 200 {object} map[string]interface{} "Pass status"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/evaluations/{evaluationId}/passed [get]
func (h *UserProgressHandler) CheckEvaluationPassed(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	evaluationID, err := strconv.ParseUint(c.Param("evaluationId"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de evaluación inválido")
		return
	}

	hasPassed, err := h.userProgressService.HasUserPassedEvaluation(uint(userID), uint(evaluationID))
	if err != nil {
		h.logger.Errorf("Error al verificar si el usuario pasó la evaluación: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo verificar el estado de la evaluación")
		return
	}

	responses.Ok(c, gin.H{"has_passed": hasPassed})
}

// @Summary Update user progress
// @Description Partially update a user progress record
// @Tags user-progress
// @Accept json
// @Produce json
// @Param id path int true "Progress ID"
// @Param data body object true "Progress update data"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-progress/{id} [patch]
func (h *UserProgressHandler) UpdateUserProgressPatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de progreso inválido")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	progress, err := h.userProgressService.UpdateUserProgressPatch(uint(id), updateData)
	if err != nil {
		h.logger.Errorf("Error al actualizar progreso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "No se pudo actualizar el progreso")
		return
	}

	responses.Ok(c, progress)
}
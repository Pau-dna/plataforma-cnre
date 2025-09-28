package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/imlargo/go-api-template/internal/dto"
	_ "github.com/imlargo/go-api-template/internal/models"

	"github.com/imlargo/go-api-template/internal/responses"
	"github.com/imlargo/go-api-template/internal/services"
)

type EnrollmentHandler struct {
	*Handler
	enrollmentService services.EnrollmentService
}

func NewEnrollmentHandler(handler *Handler, enrollmentService services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		Handler:           handler,
		enrollmentService: enrollmentService,
	}
}

// @Summary Enroll user in course
// @Description Enroll a user in a specific course
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 201 {object} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/enrollments [post]
func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
	var enrollmentData struct {
		UserID   uint `json:"user_id" binding:"required"`
		CourseID uint `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&enrollmentData); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	enrollment, err := h.enrollmentService.CreateEnrollment(enrollmentData.UserID, enrollmentData.CourseID)
	if err != nil {
		h.logger.Errorf("Error al crear la inscripción: %v", err)
		if err.Error() == "el usuario ya está inscrito en este curso" {
			responses.ErrorConflict(c, err.Error())
			return
		}
		responses.ErrorInternalServerWithMessage(c, "Error al crear la inscripción")
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

// @Summary Get enrollment by ID
// @Description Get an enrollment by its ID
// @Tags enrollments
// @Produce json
// @Param id path int true "Enrollment ID"
// @Success 200 {object} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id} [get]
func (h *EnrollmentHandler) GetEnrollment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de inscripción inválido")
		return
	}

	enrollment, err := h.enrollmentService.GetEnrollment(uint(id))
	if err != nil {
		h.logger.Errorf("Error al obtener la inscripción: %v", err)
		responses.ErrorNotFound(c, "Inscripción")
		return
	}

	responses.Ok(c, enrollment)
}

// @Summary		Update enrollment
// @Router			/api/v1/enrollments/{id} [patch]
// @Description	Update an enrollment by ID
// @Tags		enrollments
// @Param id path int true "Enrollment ID"
// @Accept		json
// @Param payload body dto.UpdateEnrollmentRequest true "Enrollment data"
// @Produce		json
// @Success		200	{object}	models.Enrollment	"Enrollment updated successfully"
// @Failure		400	{object}	responses.ErrorResponse	"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse	"Enrollment not found"
// @Failure		500	{object}	responses.ErrorResponse	"Internal Server Error"
// @Security     BearerAuth
func (h *EnrollmentHandler) UpdateEnrollmentPatch(c *gin.Context) {
	enrollmentID := c.Param("id")
	if enrollmentID == "" {
		responses.ErrorBadRequest(c, "Enrollment ID is required")
		return
	}

	enrollmentIDInt, err := strconv.Atoi(enrollmentID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Enrollment ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	enrollment, err := h.enrollmentService.UpdateEnrollmentPatch(uint(enrollmentIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, enrollment)
}

// @Summary Delete enrollment
// @Description Delete an enrollment by its ID
// @Tags enrollments
// @Param id path int true "Enrollment ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id} [delete]
func (h *EnrollmentHandler) DeleteEnrollment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de inscripción inválido")
		return
	}

	err = h.enrollmentService.DeleteEnrollment(uint(id))
	if err != nil {
		h.logger.Errorf("Error al eliminar la inscripción: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al eliminar la inscripción")
		return
	}

	responses.Ok(c, "ok")
}

// @Summary Get user enrollments
// @Description Get all enrollments for a specific user
// @Tags enrollments
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/enrollments [get]
func (h *EnrollmentHandler) GetUserEnrollments(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	enrollments, err := h.enrollmentService.GetUserEnrollments(uint(userID))
	if err != nil {
		h.logger.Errorf("Failed to get user enrollments: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al obtener la inscripcións")
		return
	}

	responses.Ok(c, enrollments)
}

// @Summary Get course enrollments
// @Description Get all enrollments for a specific course
// @Tags enrollments
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {array} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id}/enrollments [get]
func (h *EnrollmentHandler) GetCourseEnrollments(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	enrollments, err := h.enrollmentService.GetCourseEnrollments(uint(courseID))
	if err != nil {
		h.logger.Errorf("Error al obtener el curso enrollments: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al obtener la inscripcións")
		return
	}

	responses.Ok(c, enrollments)
}

// @Summary Get user course enrollment
// @Description Get enrollment for a specific user and course
// @Tags enrollments
// @Produce json
// @Param userId path int true "User ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{courseId}/enrollment [get]
func (h *EnrollmentHandler) GetUserCourseEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	enrollment, err := h.enrollmentService.GetUserCourseEnrollment(uint(userID), uint(courseID))
	if err != nil {
		h.logger.Errorf("Failed to get user course enrollment: %v", err)
		responses.ErrorNotFound(c, "Inscripción")
		return
	}

	responses.Ok(c, enrollment)
}

// @Summary Complete enrollment
// @Description Mark an enrollment as completed
// @Tags enrollments
// @Param userId path int true "User ID"
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{id}/complete [post]
func (h *EnrollmentHandler) CompleteEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	err = h.enrollmentService.CompleteEnrollment(uint(userID), uint(courseID))
	if err != nil {
		h.logger.Errorf("Failed to complete enrollment: %v", err)
		if err.Error() == "enrollment not found" {
			responses.ErrorNotFound(c, "Inscripción")
			return
		}
		responses.ErrorInternalServerWithMessage(c, "Failed to complete enrollment")
		return
	}

	responses.Ok(c, gin.H{"message": "Enrollment completed successfully"})
}

// @Summary Update enrollment progress
// @Description Update the progress of an enrollment
// @Tags enrollments
// @Accept json
// @Param userId path int true "User ID"
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{id}/progress [put]
func (h *EnrollmentHandler) UpdateProgress(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de usuario inválido")
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	var progressData struct {
		Progress float64 `json:"progress" binding:"required,min=0,max=100"`
	}

	if err := c.ShouldBindJSON(&progressData); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	err = h.enrollmentService.UpdateProgress(uint(userID), uint(courseID), progressData.Progress)
	if err != nil {
		h.logger.Errorf("Error al actualizar la inscripción progress: %v", err)
		if err.Error() == "enrollment not found" {
			responses.ErrorNotFound(c, "Inscripción")
			return
		}
		responses.ErrorInternalServerWithMessage(c, "Failed to update progress")
		return
	}

	responses.Ok(c, gin.H{"message": "Progress updated successfully"})
}

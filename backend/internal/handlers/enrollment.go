package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Param enrollment body struct{UserID uint; CourseID uint} true "Enrollment data"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrollment, err := h.enrollmentService.CreateEnrollment(enrollmentData.UserID, enrollmentData.CourseID)
	if err != nil {
		h.logger.Errorf("Failed to create enrollment: %v", err)
		if err.Error() == "user is already enrolled in this course" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enrollment"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	enrollment, err := h.enrollmentService.GetEnrollment(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get enrollment: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	c.JSON(http.StatusOK, enrollment)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	err = h.enrollmentService.DeleteEnrollment(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete enrollment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete enrollment"})
		return
	}

	c.Status(http.StatusNoContent)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	enrollments, err := h.enrollmentService.GetUserEnrollments(uint(userID))
	if err != nil {
		h.logger.Errorf("Failed to get user enrollments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

// @Summary Get course enrollments
// @Description Get all enrollments for a specific course
// @Tags enrollments
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {array} models.Enrollment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{courseId}/enrollments [get]
func (h *EnrollmentHandler) GetCourseEnrollments(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	enrollments, err := h.enrollmentService.GetCourseEnrollments(uint(courseID))
	if err != nil {
		h.logger.Errorf("Failed to get course enrollments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

// @Summary Complete enrollment
// @Description Mark an enrollment as completed
// @Tags enrollments
// @Param userId path int true "User ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{courseId}/complete [post]
func (h *EnrollmentHandler) CompleteEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	err = h.enrollmentService.CompleteEnrollment(uint(userID), uint(courseID))
	if err != nil {
		h.logger.Errorf("Failed to complete enrollment: %v", err)
		if err.Error() == "enrollment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete enrollment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment completed successfully"})
}

// @Summary Update enrollment progress
// @Description Update the progress of an enrollment
// @Tags enrollments
// @Accept json
// @Param userId path int true "User ID"
// @Param courseId path int true "Course ID"
// @Param progress body struct{Progress float64} true "Progress data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/courses/{courseId}/progress [put]
func (h *EnrollmentHandler) UpdateProgress(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var progressData struct {
		Progress float64 `json:"progress" binding:"required,min=0,max=100"`
	}

	if err := c.ShouldBindJSON(&progressData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.enrollmentService.UpdateProgress(uint(userID), uint(courseID), progressData.Progress)
	if err != nil {
		h.logger.Errorf("Failed to update enrollment progress: %v", err)
		if err.Error() == "enrollment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}

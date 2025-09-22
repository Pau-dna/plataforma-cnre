package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/internal/services"
)

type CourseHandler struct {
	*Handler
	courseService services.CourseService
}

func NewCourseHandler(handler *Handler, courseService services.CourseService) *CourseHandler {
	return &CourseHandler{
		Handler:       handler,
		courseService: courseService,
	}
}

// @Summary Create a new course
// @Description Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.Course true "Course data"
// @Success 201 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCourse, err := h.courseService.CreateCourse(&course)
	if err != nil {
		h.logger.Errorf("Failed to create course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, createdCourse)
}

// @Summary Get course by ID
// @Description Get a course by its ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.courseService.GetCourse(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get course: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Update course
// @Description Update a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param course body models.Course true "Course data"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCourse, err := h.courseService.UpdateCourse(uint(id), &course)
	if err != nil {
		h.logger.Errorf("Failed to update course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, updatedCourse)
}

// @Summary Delete course
// @Description Delete a course by its ID
// @Tags courses
// @Param id path int true "Course ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	err = h.courseService.DeleteCourse(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get all courses
// @Description Get all courses
// @Tags courses
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		h.logger.Errorf("Failed to get courses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Get course with modules
// @Description Get a course by its ID with all modules
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules [get]
func (h *CourseHandler) GetCourseWithModules(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.courseService.GetCourseWithModules(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get course with modules: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

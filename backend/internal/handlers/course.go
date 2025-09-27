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
		responses.ErrorBindJson(c, err)
		return
	}

	createdCourse, err := h.courseService.CreateCourse(&course)
	if err != nil {
		h.logger.Errorf("Error al crear el curso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al crear el curso")
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
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	course, err := h.courseService.GetCourse(uint(id))
	if err != nil {
		h.logger.Errorf("Error al obtener el curso: %v", err)
		responses.ErrorNotFound(c, "Curso")
		return
	}

	responses.Ok(c, course)
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
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	updatedCourse, err := h.courseService.UpdateCourse(uint(id), &course)
	if err != nil {
		h.logger.Errorf("Error al actualizar el curso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al actualizar el curso")
		return
	}

	responses.Ok(c, updatedCourse)
}

// @Summary		Update course
// @Router			/api/v1/courses/{id} [patch]
// @Description	Update a course by ID
// @Tags		courses
// @Param id path int true "Course ID"
// @Accept		json
// @Param payload body dto.UpdateCourseRequest true "Course data"
// @Produce		json
// @Success		200	{object}	models.Course	"Course updated successfully"
// @Failure		400	{object}	responses.ErrorResponse	"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse	"Course not found"
// @Failure		500	{object}	responses.ErrorResponse	"Internal Server Error"
// @Security     BearerAuth
func (h *CourseHandler) UpdateCoursePatch(c *gin.Context) {
	courseID := c.Param("id")
	if courseID == "" {
		responses.ErrorBadRequest(c, "Course ID is required")
		return
	}

	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Course ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	course, err := h.courseService.UpdateCoursePatch(uint(courseIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, course)
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
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	err = h.courseService.DeleteCourse(uint(id))
	if err != nil {
		h.logger.Errorf("Error al eliminar el curso: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al eliminar el curso")
		return
	}

	responses.Ok(c, "ok")
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
		h.logger.Errorf("Error al obtener el cursos: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Error al obtener el cursos")
		return
	}

	responses.Ok(c, courses)
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
		responses.ErrorBadRequest(c, "ID de curso inválido")
		return
	}

	course, err := h.courseService.GetCourseWithModules(uint(id))
	if err != nil {
		h.logger.Errorf("Error al obtener el curso with modules: %v", err)
		responses.ErrorNotFound(c, "Curso")
		return
	}

	responses.Ok(c, course)
}

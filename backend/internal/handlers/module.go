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

type ModuleHandler struct {
	*Handler
	moduleService services.ModuleService
}

func NewModuleHandler(handler *Handler, moduleService services.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		Handler:       handler,
		moduleService: moduleService,
	}
}

// @Summary Create a new module
// @Description Create a new module
// @Tags modules
// @Accept json
// @Produce json
// @Param module body models.Module true "Module data"
// @Success 201 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/modules [post]
func (h *ModuleHandler) CreateModule(c *gin.Context) {
	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	createdModule, err := h.moduleService.CreateModule(&module)
	if err != nil {
		h.logger.Errorf("Failed to create module: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to create module")
		return
	}

	c.JSON(http.StatusCreated, createdModule)
}

// @Summary Get module by ID
// @Description Get a module by its ID
// @Tags modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/modules/{id} [get]
func (h *ModuleHandler) GetModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid module ID")
		return
	}

	module, err := h.moduleService.GetModule(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get module: %v", err)
		responses.ErrorNotFound(c, "Module")
		return
	}

	responses.Ok(c, module)
}

// @Summary Update module
// @Description Update a module by its ID
// @Tags modules
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param module body models.Module true "Module data"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/modules/{id} [put]
func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid module ID")
		return
	}

	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	updatedModule, err := h.moduleService.UpdateModule(uint(id), &module)
	if err != nil {
		h.logger.Errorf("Failed to update module: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to update module")
		return
	}

	responses.Ok(c, updatedModule)
}

// @SummaryUpdate module
// @Router/api/v1/modules/{id} [patch]
// @DescriptionUpdate a module by ID
// @Tagsmodules
// @Param id path int true "Module ID"
// @Acceptjson
// @Param payload body dto.UpdateModuleRequest true "Module data"
// @Producejson
// @Success200{object}models.Module"Module updated successfully"
// @Failure400{object}responses.ErrorResponse"Bad Request"
// @Failure404{object}responses.ErrorResponse"Module not found"
// @Failure500{object}responses.ErrorResponse"Internal Server Error"
// @Security     BearerAuth
func (h *ModuleHandler) UpdateModulePatch(c *gin.Context) {
	moduleID := c.Param("id")
	if moduleID == "" {
		responses.ErrorBadRequest(c, "Module ID is required")
		return
	}

	moduleIDInt, err := strconv.Atoi(moduleID)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid Module ID: "+err.Error())
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	module, err := h.moduleService.UpdateModulePatch(uint(moduleIDInt), payload)
	if err != nil {
		responses.ErrorInternalServerWithMessage(c, err.Error())
		return
	}

	responses.Ok(c, module)
}

// @Summary Delete module
// @Description Delete a module by its ID
// @Tags modules
// @Param id path int true "Module ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/modules/{id} [delete]
func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid module ID")
		return
	}

	err = h.moduleService.DeleteModule(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete module: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to delete module")
		return
	}

	responses.Ok(c, "ok")
}

// @Summary Get modules by course
// @Description Get all modules for a specific course
// @Tags modules
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {array} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{courseId}/modules [get]
func (h *ModuleHandler) GetModulesByCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid course ID")
		return
	}

	modules, err := h.moduleService.GetModulesByCourse(uint(courseID))
	if err != nil {
		h.logger.Errorf("Failed to get modules by course: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to get modules")
		return
	}

	responses.Ok(c, modules)
}

// @Summary Get module with content
// @Description Get a module by its ID with all content
// @Tags modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/modules/{id}/content [get]
func (h *ModuleHandler) GetModuleWithContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid module ID")
		return
	}

	module, err := h.moduleService.GetModuleWithContent(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get module with content: %v", err)
		responses.ErrorNotFound(c, "Module")
		return
	}

	responses.Ok(c, module)
}

// @Summary Reorder modules
// @Description Reorder modules for a course
// @Tags modules
// @Accept json
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{courseId}/modules/reorder [post]
func (h *ModuleHandler) ReorderModules(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		responses.ErrorBadRequest(c, "Invalid course ID")
		return
	}

	var moduleOrders []struct {
		ID    uint `json:"id" binding:"required"`
		Order int  `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&moduleOrders); err != nil {
		responses.ErrorBindJson(c, err)
		return
	}

	// Convert to the expected type
	var convertedOrders []struct {
		ID    uint
		Order int
	}
	for _, order := range moduleOrders {
		convertedOrders = append(convertedOrders, struct {
			ID    uint
			Order int
		}{
			ID:    order.ID,
			Order: order.Order,
		})
	}

	err = h.moduleService.ReorderModules(uint(courseID), convertedOrders)
	if err != nil {
		h.logger.Errorf("Failed to reorder modules: %v", err)
		responses.ErrorInternalServerWithMessage(c, "Failed to reorder modules")
		return
	}

	responses.Ok(c, gin.H{"message": "Modules reordered successfully"})
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/models"
"github.com/imlargo/go-api-template/internal/responses"
	"github.com/imlargo/go-api-template/internal/services"
)

type ContentHandler struct {
	*Handler
	contentService services.ContentService
}

func NewContentHandler(handler *Handler, contentService services.ContentService) *ContentHandler {
	return &ContentHandler{
		Handler:        handler,
		contentService: contentService,
	}
}

// @Summary Create content
// @Description Create new content for a module
// @Tags content
// @Accept json
// @Produce json
// @Param content body models.Content true "Content data"
// @Success 201 {object} models.Content
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/content [post]
func (h *ContentHandler) CreateContent(c *gin.Context) {
	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdContent, err := h.contentService.CreateContent(&content)
	if err != nil {
		h.logger.Errorf("Failed to create content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
		return
	}

	c.JSON(http.StatusCreated, createdContent)
}

// @Summary Get content
// @Description Get content by ID
// @Tags content
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} models.Content
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/content/{id} [get]
func (h *ContentHandler) GetContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	content, err := h.contentService.GetContent(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get content: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// @Summary Update content
// @Description Update content by ID
// @Tags content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param content body models.Content true "Content data"
// @Success 200 {object} models.Content
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/content/{id} [put]
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedContent, err := h.contentService.UpdateContent(uint(id), &content)
	if err != nil {
		h.logger.Errorf("Failed to update content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update content"})
		return
	}

	c.JSON(http.StatusOK, updatedContent)
}

// @SummaryUpdate content
// @Router/api/v1/contents/{id} [patch]
// @DescriptionUpdate a content by ID
// @Tagscontents
// @Param id path int true "Content ID"
// @Acceptjson
// @Param payload body dto.UpdateContentRequest true "Content data"
// @Producejson
// @Success200{object}models.Content"Content updated successfully"
// @Failure400{object}responses.ErrorResponse"Bad Request"
// @Failure404{object}responses.ErrorResponse"Content not found"
// @Failure500{object}responses.ErrorResponse"Internal Server Error"
// @Security     BearerAuth
func (h *ContentHandler) UpdateContentPatch(c *gin.Context) {
contentID := c.Param("id")
if contentID == "" {
responses.ErrorBadRequest(c, "Content ID is required")
return
}

contentIDInt, err := strconv.Atoi(contentID)
if err != nil {
responses.ErrorBadRequest(c, "Invalid Content ID: "+err.Error())
return
}

var payload map[string]interface{}
if err := c.BindJSON(&payload); err != nil {
responses.ErrorBadRequest(c, "Invalid request payload: "+err.Error())
return
}

content, err := h.contentService.UpdateContentPatch(uint(contentIDInt), payload)
if err != nil {
responses.ErrorInternalServerWithMessage(c, err.Error())
return
}

responses.Ok(c, content)
}

// @Summary Delete content
// @Description Delete content by ID
// @Tags content
// @Param id path int true "Content ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/content/{id} [delete]
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	err = h.contentService.DeleteContent(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to delete content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete content"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get module contents
// @Description Get all content for a specific module
// @Tags content
// @Produce json
// @Param moduleId path int true "Module ID"
// @Success 200 {array} models.Content
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/modules/{moduleId}/content [get]
func (h *ContentHandler) GetContentsByModule(c *gin.Context) {
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	contents, err := h.contentService.GetContentsByModule(uint(moduleID))
	if err != nil {
		h.logger.Errorf("Failed to get contents by module: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get contents"})
		return
	}

	c.JSON(http.StatusOK, contents)
}

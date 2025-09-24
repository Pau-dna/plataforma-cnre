package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func ErrorBindJson(c *gin.Context, err error) {
	NewErrorResponse(c, http.StatusBadRequest, err.Error(), errBindJson, nil)
}

func ErrorNotFound(c *gin.Context, model string) {
	NewErrorResponse(c, http.StatusNotFound, model+" not found", errNotFound, nil)
}

func ErrorInternalServer(c *gin.Context) {
	NewErrorResponse(c, http.StatusInternalServerError, "internal server error", errInternalServer, nil)
}

func ErrorInternalServerWithMessage(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusInternalServerError, message, errInternalServer, nil)
}

func ErrorBadRequest(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusBadRequest, message, errBadRequest, nil)
}

func ErrorTooManyRequests(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusTooManyRequests, message, errTooManyRequests, nil)
}

func ErrorUnauthorized(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusUnauthorized, message, errUnauthorized, nil)
}

func ErrorConflict(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusConflict, message, errConflict, nil)
}

func NewErrorResponse(c *gin.Context, httpStatusCode int, message string, code string, payload map[string]interface{}) {
	c.JSON(httpStatusCode, ErrorResponse{
		Code:    code,
		Message: message,
		Payload: payload,
	})
}

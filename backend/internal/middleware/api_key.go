package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/responses"
)

func ApiKeyMiddleware(apiKey string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		apiKeyHeader := ctx.GetHeader("X-API-Key")

		if apiKeyHeader == "" {
			ctx.Abort()
			responses.ErrorUnauthorized(ctx, "falta el header de autorización")
			return
		}

		if apiKeyHeader != apiKey {
			ctx.Abort()
			responses.ErrorUnauthorized(ctx, "clave de API inválida")
			return
		}

		ctx.Next()
	}
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/responses"
)

func BearerApiKeyMiddleware(apiKey string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.Abort()
			responses.ErrorUnauthorized(ctx, "falta el header de autorización")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.Abort()
			responses.ErrorUnauthorized(ctx, "el header de autorización debe tener el formato 'Bearer token'")
			return
		}

		apiKeyHeader := parts[1]
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

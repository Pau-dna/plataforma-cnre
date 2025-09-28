package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/internal/responses"
	"github.com/imlargo/go-api-template/pkg/ratelimiter"
)

func NewRateLimiterMiddleware(rl ratelimiter.RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		allow, retryAfter := rl.Allow(ip)
		if !allow {
			message := "Límite de solicitudes excedido. Intenta de nuevo en " + fmt.Sprintf("%.2f", retryAfter) + " segundos"
			responses.ErrorTooManyRequests(ctx, message)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

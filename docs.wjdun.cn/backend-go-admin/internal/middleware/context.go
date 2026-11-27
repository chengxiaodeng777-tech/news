package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

const AppContextKey = "app_ctx"

func AppContext(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AppContextKey, ctx)
		c.Next()
	}
}

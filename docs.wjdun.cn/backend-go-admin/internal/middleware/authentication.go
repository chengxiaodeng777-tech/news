package middleware

import (
	"net/http"

	authservice "base/internal/service/auth"

	"github.com/gin-gonic/gin"
)

const (
	CurrentUserIDKey   = "current_user_id"
	CurrentUsernameKey = "current_username"
	CurrentRolesKey    = "current_roles"
)

func Authentication(service *authservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := authservice.NormalizeToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			return
		}

		claims, err := service.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效"})
			return
		}

		c.Set(CurrentUserIDKey, claims.UserID)
		c.Set(CurrentUsernameKey, claims.Username)
		c.Set(CurrentRolesKey, claims.Roles)
		c.Next()
	}
}

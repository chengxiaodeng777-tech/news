package middleware

import (
	"net/http"

	"base/ent"
	authservice "base/internal/service/auth"

	casbinlib "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Permission reads the current role-menu assignments. The enforcer passed by
// the existing router is a snapshot and can be stale after permission changes.
func Permission(enforcer *casbinlib.Enforcer, permissionCode string) gin.HandlerFunc {
	return PermissionAny(enforcer, permissionCode)
}

func PermissionAny(_ *casbinlib.Enforcer, permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get(CurrentUserIDKey)
		userID, ok := userIDValue.(int)
		if !exists || !ok || userID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "当前用户信息无效"})
			return
		}

		databaseValue, exists := c.Get(DatabaseKey)
		client, ok := databaseValue.(*ent.Client)
		if !exists || !ok || client == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "数据库不可用"})
			return
		}

		allowed, err := authservice.New(client, "").HasAnyPermission(c.Request.Context(), userID, permissionCodes...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限校验失败"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无访问权限"})
			return
		}
		c.Next()
	}
}

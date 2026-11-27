package middleware

import (
	"net/http"

	casbinlib "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Permission(enforcer *casbinlib.Enforcer, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameValue, exists := c.Get(CurrentUsernameKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未获取到当前用户"})
			return
		}
		username, ok := usernameValue.(string)
		if !ok || username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "当前用户信息无效"})
			return
		}

		rolesValue, _ := c.Get(CurrentRolesKey)
		roleList, _ := rolesValue.([]string)

		for _, roleCode := range roleList {
			if roleCode == "admin" {
				c.Next()
				return
			}
		}

		subject := "user:" + username
		ok, err := enforcer.Enforce(subject, permissionCode, "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限校验失败"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无访问权限"})
			return
		}
		c.Next()
	}
}

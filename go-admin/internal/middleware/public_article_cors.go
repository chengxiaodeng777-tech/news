package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Public article reads may be embedded by independently hosted websites.
func PublicArticleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if (path == "/api/articles" || strings.HasPrefix(path, "/api/articles/")) &&
			(c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		}
		c.Next()
	}
}

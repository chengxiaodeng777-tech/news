package middleware

import (
	"bytes"
	"io"

	"base/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Audit(log *zap.Logger, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		requestID, _ := c.Get(RequestIDKey)
		username, _ := c.Get(CurrentUsernameKey)
		log.Info(
			"审计日志",
			logger.String("request_id", toString(requestID)),
			logger.String("username", toString(username)),
			logger.String("action", action),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.String("request_body", requestBody),
			zap.Int("status", c.Writer.Status()),
		)
	}
}

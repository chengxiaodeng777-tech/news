package middleware

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	DatabaseKey = "database"
	RedisKey    = "redis"
	LoggerKey   = "logger"
)

func Resources(database any, redis any, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DatabaseKey, database)
		c.Set(RedisKey, redis)
		c.Set(LoggerKey, log)
		c.Next()
	}
}

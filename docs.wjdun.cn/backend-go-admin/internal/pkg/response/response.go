package response

import (
	"net/http"

	"base/internal/logger"
	"base/internal/pkg/apperror"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

func AppError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if logValue, exists := c.Get("logger"); exists {
		requestID, _ := c.Get("request_id")
		if log, ok := logValue.(interface {
			Error(string, ...logger.Field)
		}); ok {
			log.Error(
				"请求处理失败",
				logger.String("request_id", toString(requestID)),
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.Error(err),
			)
		}
	}

	appErr, ok := err.(*apperror.Error)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": apperror.CodeInternalServer,
			"msg":  "系统内部错误",
			"data": nil,
		})
		return
	}

	status := http.StatusInternalServerError
	switch {
	case appErr.Code >= 40000 && appErr.Code < 40100:
		status = http.StatusBadRequest
	case appErr.Code >= 40100 && appErr.Code < 40300:
		status = http.StatusUnauthorized
	case appErr.Code >= 40300 && appErr.Code < 40400:
		status = http.StatusForbidden
	case appErr.Code >= 40400 && appErr.Code < 40900:
		status = http.StatusNotFound
	case appErr.Code >= 40900 && appErr.Code < 50000:
		status = http.StatusConflict
	}

	c.JSON(status, gin.H{
		"code": appErr.Code,
		"msg":  appErr.Message,
		"data": nil,
	})
}

func toString(v any) string {
	s, _ := v.(string)
	return s
}

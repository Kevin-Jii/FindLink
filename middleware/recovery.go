package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"app/common"
	"app/utils/logger"
)

// RecoveryMiddleware 统一panic恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				traceID := GetTraceID(c)
				logger.Error("panic recovered",
					logger.String("trace_id", traceID),
					logger.String("path", c.Request.URL.Path),
					logger.Any("error", err),
					logger.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":     common.ServerErr.Code,
					"msg":      common.ServerErr.Msg,
					"trace_id": traceID,
				})
			}
		}()
		c.Next()
	}
}

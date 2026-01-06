package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDKey = "X-Trace-ID"

// TraceMiddleware 请求追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从请求头获取，没有则生成新的
		traceID := c.GetHeader(TraceIDKey)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 设置到上下文和响应头
		c.Set(TraceIDKey, traceID)
		c.Header(TraceIDKey, traceID)

		c.Next()
	}
}

// GetTraceID 从上下文获取TraceID
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		return traceID.(string)
	}
	return ""
}

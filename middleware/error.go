package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/common"
	"app/utils/logger"
)

// ErrorResponse 统一错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	TraceID string `json:"trace_id,omitempty"`
}

// ErrorMiddleware 统一错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			traceID := GetTraceID(c)

			logger.Error("request error",
				logger.String("trace_id", traceID),
				logger.String("path", c.Request.URL.Path),
				logger.String("error", err.Error()),
			)

			// 根据错误类型返回不同状态码
			var errno common.Errno
			switch e := err.Err.(type) {
			case common.Errno:
				errno = e
			default:
				errno = common.ServerErr
			}

			httpStatus := http.StatusOK
			if errno.Code == 401 {
				httpStatus = http.StatusUnauthorized
			} else if errno.Code == 403 {
				httpStatus = http.StatusForbidden
			}

			c.JSON(httpStatus, ErrorResponse{
				Code:    errno.Code,
				Msg:     errno.Msg,
				TraceID: traceID,
			})
		}
	}
}

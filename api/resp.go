package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/common"
	"app/consts"
	"app/middleware"
)

type Resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id,omitempty"`
}

// WriteResp 统一响应输出
func WriteResp(ctx *gin.Context, data any, errno common.Errno) {
	traceID := middleware.GetTraceID(ctx)
	ctx.JSON(http.StatusOK, Resp{
		Code:    errno.Code,
		Msg:     errno.Msg,
		Data:    data,
		TraceID: traceID,
	})
}

// WriteSuccess 成功响应
func WriteSuccess(ctx *gin.Context, data any) {
	WriteResp(ctx, data, common.OK)
}

// WriteError 错误响应
func WriteError(ctx *gin.Context, errno common.Errno) {
	WriteResp(ctx, nil, errno)
}

func GetUserFromCtx(ctx *gin.Context) *common.User {
	user, exist := ctx.Get(consts.CustomerUserKey)
	if !exist {
		return nil
	}
	return user.(*common.User)
}

func GetAdminUserFromCtx(ctx *gin.Context) *common.AdminUser {
	user, exist := ctx.Get(consts.AdminUserKey)
	if !exist {
		return nil
	}
	return user.(*common.AdminUser)
}

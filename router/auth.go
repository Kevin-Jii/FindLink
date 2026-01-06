package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/common"
	"app/consts"
)

type TokenFun func(ctx context.Context, token string) (*common.User, error)
type TokenAdminFun func(ctx context.Context, token string) (*common.AdminUser, error)

// AuthMiddleware C端用户鉴权中间件
func AuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 白名单放行
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(consts.UserTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		user, err := getTokenFun(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		if user == nil || user.UserID == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		ctx.Set(consts.CustomerUserKey, user)
		ctx.Next()
	}
}

// AdminAuthMiddleware 管理后台鉴权中间件
func AdminAuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenAdminFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 白名单放行
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(consts.AdminTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		user, err := getTokenFun(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		if user == nil || user.UserID == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		ctx.Set(consts.AdminUserKey, user)
		ctx.Next()
	}
}

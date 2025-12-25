package admin

import (
	"app/api"
	"app/common"
	"app/service/dto"

	"github.com/gin-gonic/gin"
)

// GetUserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  api.Resp{data=dto.UserInfoResp}
// @Router       /admin/v1/user/info [get]
func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body  dto.CreateUserReq  true  "用户信息"
// @Success      200  {object}  api.Resp{data=object{id=int64}}
// @Router       /admin/v1/user/create [post]
func (c *Ctrl) CreateUser(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.CreateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}
	userId, errno := c.user.CreateUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]int64{
		"id": userId,
	}, errno)
}

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body  dto.UpdateUserReq  true  "用户信息"
// @Success      200  {object}  api.Resp
// @Router       /admin/v1/user/update [post]
func (c *Ctrl) UpdateUser(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}
	errno := c.user.UpdateUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// UpdateUserStatus 更新用户状态
// @Summary      更新用户状态
// @Description  启用或禁用用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body  dto.UpdateUserStatusReq  true  "状态信息"
// @Success      200  {object}  api.Resp
// @Router       /admin/v1/user/status [post]
func (c *Ctrl) UpdateUserStatus(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}
	errno := c.user.UpdateUserStatus(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

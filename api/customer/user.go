package customer

import (
	"app/api"
	"app/common"
	"app/service/dto"

	"github.com/gin-gonic/gin"
)

// Login C端用户登录
// @Summary      C端用户登录
// @Description  使用手机号密码登录
// @Tags         C端-用户
// @Accept       json
// @Produce      json
// @Param        request  body  dto.UserLoginReq  true  "登录参数"
// @Success      200  {object}  api.Resp{data=dto.UserLoginResp}
// @Router       /customer/v1/user/login [post]
func (c *Ctrl) Login(ctx *gin.Context) {
	req := &dto.UserLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.User.Login(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// Register C端用户注册
// @Summary      C端用户注册
// @Description  使用手机号注册
// @Tags         C端-用户
// @Accept       json
// @Produce      json
// @Param        request  body  dto.UserRegisterReq  true  "注册参数"
// @Success      200  {object}  api.Resp{data=dto.UserLoginResp}
// @Router       /customer/v1/user/register [post]
func (c *Ctrl) Register(ctx *gin.Context) {
	req := &dto.UserRegisterReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.User.Register(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// GetUserInfo C端获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户信息
// @Tags         C端-用户
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  api.Resp{data=dto.CustomerUserInfoResp}
// @Router       /customer/v1/user/info [get]
func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.User.GetUserInfo(ctx.Request.Context(), user.UserID)
	api.WriteResp(ctx, resp, errno)
}

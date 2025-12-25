package admin

import (
	"app/api"
	"app/common"
	"app/service/dto"

	"github.com/gin-gonic/gin"
)

// GetSmsCodeCaptcha 获取滑块验证码
// @Summary      获取滑块验证码
// @Description  获取滑块验证码图片
// @Tags         验证码
// @Accept       json
// @Produce      json
// @Param        once  query  string  false  "随机字符串"
// @Param        ts    query  int     false  "时间戳"
// @Param        sign  query  string  false  "签名"
// @Success      200  {object}  api.Resp{data=dto.GetVerifyCaptchaResp}
// @Router       /admin/v1/user/verify/captcha [get]
func (c *Ctrl) GetSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.GetVerifyCaptchaReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.GetSlideCaptcha(ctx.Request.Context())
	api.WriteResp(ctx, resp, errno)
}

// CheckSmsCodeCaptcha 校验滑块验证码
// @Summary      校验滑块验证码
// @Description  校验滑块验证码位置
// @Tags         验证码
// @Accept       json
// @Produce      json
// @Param        request  body  dto.CheckCaptchaReq  true  "验证码校验参数"
// @Success      200  {object}  api.Resp{data=dto.CheckCaptchaDtoResp}
// @Router       /admin/v1/user/verify/captcha/check [post]
func (c *Ctrl) CheckSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.CheckCaptchaReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.CheckSlideCaptcha(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// Login 用户登录
// @Summary      用户登录
// @Description  使用用户名密码登录
// @Tags         登录
// @Accept       json
// @Produce      json
// @Param        request  body  dto.LoginReq  true  "登录参数"
// @Success      200  {object}  api.Resp{data=dto.LoginResp}
// @Router       /admin/v1/user/login [post]
func (c *Ctrl) Login(ctx *gin.Context) {
	req := &dto.LoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.Login(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

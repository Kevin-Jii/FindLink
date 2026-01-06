package dto

import (
	"app/utils/tools"
	"fmt"
)

type GetVerifyCaptchaReq struct {
	Once string `form:"once"`
	Time int64  `form:"ts"`
	Sign string `form:"sign"` // 秘钥固定加密： md5(once+daqing2025+ts) 转小写
}

func (r *GetVerifyCaptchaReq) CheckSign() bool {
	return r.Sign == tools.Sha256Hash(fmt.Sprintf("%s%s%d", r.Once, "daqing2025", r.Time))
}

type GetVerifyCaptchaResp struct {
	Key            string `json:"key"`
	ImageBs64      string `json:"image_base64"`       // 包含“data:image/jpeg;base64
	TitleImageBs64 string `json:"title_image_base64"` // 滑块图片，包含“data:image/jpeg;base64
	TitleHeight    int    `json:"title_height"`       // 滑块图片高
	TitleWidth     int    `json:"title_width"`        // 滑块图片宽
	TitleX         int    `json:"title_x"`            // 滑块图的x坐标
	TitleY         int    `json:"title_y"`            // 滑块图的y坐标
	Expire         int64  `json:"expire"`             // 过期时间
}

type CheckCaptchaReq struct {
	Key    string `json:"key"`
	SlideX int    `json:"slide_x"`
	SlideY int    `json:"slide_y"`
}

type CheckCaptchaDtoResp struct {
	Ticket string `json:"ticket"`
	Expire int64  `json:"expire"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required,min=2,max=32"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=32"` // 密码
	Ticket   string `json:"ticket"`                                   // 验证码ticket（可选）
}

// LoginResp 登录响应
type LoginResp struct {
	Token    string `json:"token"`     // 访问令牌
	ExpireAt int64  `json:"expire_at"` // 过期时间戳
	UserID   int64  `json:"user_id"`   // 用户ID
	Name     string `json:"name"`      // 用户名
}

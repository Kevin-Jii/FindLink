package dto

import "app/adaptor/repo/model"

// UserLoginReq C端用户登录请求
type UserLoginReq struct {
	Mobile   string `json:"mobile" binding:"required"`   // 手机号
	Password string `json:"password" binding:"required"` // 密码
}

// UserLoginResp C端用户登录响应
type UserLoginResp struct {
	Token    string `json:"token"`     // 访问令牌
	ExpireAt int64  `json:"expire_at"` // 过期时间戳
	UserID   int64  `json:"user_id"`   // 用户ID
	Nickname string `json:"nickname"`  // 昵称
}

// UserRegisterReq C端用户注册请求
type UserRegisterReq struct {
	Mobile   string `json:"mobile" binding:"required"`   // 手机号
	Password string `json:"password" binding:"required"` // 密码
	Nickname string `json:"nickname"`                    // 昵称
}

// CustomerUserInfoResp C端用户信息响应
type CustomerUserInfoResp struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
	Gender   int32  `json:"gender"`
}

// UserModel 用户模型转换
type UserModel struct {
	ID       int64
	Nickname string
	Mobile   string
	Password string
	Status   int32
}

func (u *UserModel) ToModel() *model.User {
	return &model.User{
		ID:       u.ID,
		Nickname: u.Nickname,
		Mobile:   u.Mobile,
		Password: u.Password,
		Status:   u.Status,
	}
}

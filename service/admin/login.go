package admin

import (
	"context"
	"time"

	"app/common"
	"app/service/dto"
	"app/utils/tools"
)

// Login 管理员登录
func (s *Service) Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginResp, common.Errno) {
	user, err := s.adminUser.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, common.UserNotFoundErr
	}

	// 使用bcrypt验证密码
	if !tools.CheckPassword(req.Password, user.Password) {
		return nil, common.AuthErr.WithMsg("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, common.AuthErr.WithMsg("账户已被禁用")
	}

	// 生成token，添加admin前缀区分
	token := tools.UUIDHex()
	expireAt := time.Now().Add(time.Hour * 24).Unix()

	// 存储token到redis
	err = s.verify.SetToken(ctx, "admin:"+token, user.ID, time.Hour*24)
	if err != nil {
		return nil, common.RedisErr
	}

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: expireAt,
		UserID:   user.ID,
		Name:     user.Name,
	}, common.OK
}

package admin

import (
	"context"
	"time"

	"app/common"
	"app/service/dto"
	"app/utils/tools"
)

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginResp, common.Errno) {
	// 查询用户（这里简化处理，实际应该查数据库）
	user, err := s.adminUser.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, common.UserNotFoundErr
	}

	// 验证密码（实际应该用加密比对）
	if user.Password != tools.Sha256Hash(req.Password) {
		return nil, common.AuthErr.WithMsg("密码错误")
	}

	// 生成token
	token := tools.UUIDHex()
	expireAt := time.Now().Add(time.Hour * 24).Unix()

	// 存储token到redis
	err = s.verify.SetToken(ctx, token, user.ID, time.Hour*24)
	if err != nil {
		return nil, common.RedisErr.WithErr(err)
	}

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: expireAt,
		UserID:   user.ID,
		Name:     user.Name,
	}, common.OK
}

package user

import (
	"context"

	"app/adaptor"
	"app/adaptor/redis"
	"app/adaptor/repo/user"
)

type Service struct {
	userRepo user.IUser
	verify   redis.IVerify
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		userRepo: user.NewUser(adaptor),
		verify:   redis.NewVerify(adaptor.GetRedis()),
	}
}

// Logout 用户退出登录
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.verify.DelToken(ctx, token)
}

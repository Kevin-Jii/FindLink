package user

import (
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
		verify:   redis.NewVerify(adaptor),
	}
}

package admin

import (
	"github.com/wenlng/go-captcha/v2/slide"
	"app/adaptor"
	"app/adaptor/redis"
	"app/adaptor/repo/admin"
	"app/utils/captcha"
)

type Service struct {
	adminUser admin.IAdminUser
	verify    redis.IVerify
	captcha   slide.Captcha
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUser(adaptor),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
	}
}

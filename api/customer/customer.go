package customer

import (
	"app/adaptor"
	"app/service/user"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *user.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    user.NewService(adaptor),
	}
}

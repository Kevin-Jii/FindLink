package user

import (
	"context"

	"github.com/go-redis/redis"
	"gorm.io/gorm"

	"app/adaptor"
	"app/adaptor/repo/model"
)

type IUser interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByMobile(ctx context.Context, mobile string) (*model.User, error)
	GetByOpenID(ctx context.Context, openID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateStatus(ctx context.Context, id int64, status int32) error
}

type User struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUser(adaptor adaptor.IAdaptor) *User {
	return &User{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (u *User) GetByMobile(ctx context.Context, mobile string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	return &user, err
}

func (u *User) GetByOpenID(ctx context.Context, openID string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("open_id = ?", openID).First(&user).Error
	return &user, err
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Save(user).Error
}

func (u *User) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

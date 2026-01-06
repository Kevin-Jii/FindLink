package user

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"app/adaptor/repo/model"
	"app/common"
	"app/service/dto"
	"app/utils/tools"
)

// Login C端用户登录
func (s *Service) Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserLoginResp, common.Errno) {
	user, err := s.userRepo.GetByMobile(ctx, req.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		return nil, common.DatabaseErr
	}

	// 使用bcrypt验证密码
	if !tools.CheckPassword(req.Password, user.Password) {
		return nil, common.AuthErr.WithMsg("密码错误")
	}

	if user.Status != 1 {
		return nil, common.AuthErr.WithMsg("账户已被禁用")
	}

	token := tools.UUIDHex()
	expireAt := time.Now().Add(time.Hour * 24 * 7).Unix()

	err = s.verify.SetToken(ctx, "user:"+token, user.ID, time.Hour*24*7)
	if err != nil {
		return nil, common.RedisErr
	}

	return &dto.UserLoginResp{
		Token:    token,
		ExpireAt: expireAt,
		UserID:   user.ID,
		Nickname: user.Nickname,
	}, common.OK
}

// Register C端用户注册
func (s *Service) Register(ctx context.Context, req *dto.UserRegisterReq) (*dto.UserLoginResp, common.Errno) {
	// 检查手机号是否已注册
	_, err := s.userRepo.GetByMobile(ctx, req.Mobile)
	if err == nil {
		return nil, common.ParamErr.WithMsg("手机号已注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.DatabaseErr
	}

	// 使用bcrypt加密密码
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return nil, common.ServerErr
	}

	// 创建用户
	userModel := &model.User{
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Password: hashedPassword,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, userModel); err != nil {
		return nil, common.DatabaseErr
	}

	// 自动登录
	token := tools.UUIDHex()
	expireAt := time.Now().Add(time.Hour * 24 * 7).Unix()

	_ = s.verify.SetToken(ctx, "user:"+token, userModel.ID, time.Hour*24*7)

	return &dto.UserLoginResp{
		Token:    token,
		ExpireAt: expireAt,
		UserID:   userModel.ID,
		Nickname: userModel.Nickname,
	}, common.OK
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (*dto.CustomerUserInfoResp, common.Errno) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		return nil, common.DatabaseErr
	}

	return &dto.CustomerUserInfoResp{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
	}, common.OK
}

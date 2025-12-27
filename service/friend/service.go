package friend

import (
	"context"
	"fmt"

	"app/adaptor/repo/friend"
	"app/adaptor/repo/model"
	"app/common"
	"app/service/dto"
)

// IFriendService 好友服务接口
type IFriendService interface {
	SendFriendRequest(ctx context.Context, fromUserID, toUserID int64, message string) error
	GetFriendRequests(ctx context.Context, userID int64) ([]*dto.FriendRequestResp, error)
	AcceptFriendRequest(ctx context.Context, userID, requestID int64) error
	RejectFriendRequest(ctx context.Context, userID, requestID int64) error
	GetFriendList(ctx context.Context, userID int64) ([]*dto.FriendResp, error)
	RemoveFriend(ctx context.Context, userID, friendID int64) error
	SearchUsers(ctx context.Context, userID int64, keyword string) ([]*dto.UserSearchResp, error)
}

// FriendService 好友服务实现
type FriendService struct {
	repo *friend.FriendRepository
}

// NewFriendService 创建好友服务
func NewFriendService(repo *friend.FriendRepository) *FriendService {
	return &FriendService{repo: repo}
}

// SendFriendRequest 发送好友请求
func (s *FriendService) SendFriendRequest(ctx context.Context, fromUserID, toUserID int64, message string) error {
	if fromUserID == toUserID {
		return common.CannotAddSelfErr
	}

	// 检查是否已经是好友
	areFriends, err := s.repo.AreFriends(ctx, fromUserID, toUserID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if areFriends {
		return common.AlreadyFriendsErr
	}

	// 检查是否已有待处理的请求
	requests, err := s.repo.GetPendingRequests(ctx, toUserID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	for _, req := range requests {
		if req.FromUserID == fromUserID {
			return common.FriendRequestExistsErr
		}
	}

	req := &model.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     model.FriendStatusPending,
	}

	return s.repo.CreateFriendRequest(ctx, req)
}

// GetFriendRequests 获取好友请求列表
func (s *FriendService) GetFriendRequests(ctx context.Context, userID int64) ([]*dto.FriendRequestResp, error) {
	requests, err := s.repo.GetPendingRequests(ctx, userID)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}

	resp := make([]*dto.FriendRequestResp, len(requests))
	for i, req := range requests {
		resp[i] = &dto.FriendRequestResp{
			ID:        req.ID,
			FromUserID: req.FromUserID,
			Message:   req.Message,
			Status:    string(req.Status),
			CreatedAt: req.CreatedAt,
			// TODO: 获取用户信息
			Nickname: fmt.Sprintf("用户%d", req.FromUserID),
			Avatar:   "",
		}
	}

	return resp, nil
}

// AcceptFriendRequest 接受好友请求
func (s *FriendService) AcceptFriendRequest(ctx context.Context, userID, requestID int64) error {
	req, err := s.repo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if req == nil {
		return common.FriendNotFoundErr
	}
	if req.ToUserID != userID {
		return common.PermissionErr
	}

	return s.repo.AcceptFriendRequest(ctx, requestID)
}

// RejectFriendRequest 拒绝好友请求
func (s *FriendService) RejectFriendRequest(ctx context.Context, userID, requestID int64) error {
	req, err := s.repo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if req == nil {
		return common.FriendNotFoundErr
	}
	if req.ToUserID != userID {
		return common.PermissionErr
	}

	return s.repo.RejectFriendRequest(ctx, requestID)
}

// GetFriendList 获取好友列表
func (s *FriendService) GetFriendList(ctx context.Context, userID int64) ([]*dto.FriendResp, error) {
	friends, err := s.repo.GetFriends(ctx, userID, string(model.FriendStatusAccepted))
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}

	resp := make([]*dto.FriendResp, len(friends))
	for i, f := range friends {
		resp[i] = &dto.FriendResp{
			ID:            f.ID,
			UserID:        f.UserID,
			FriendID:      f.FriendID,
			SharingStatus: string(f.SharingStatus),
			// TODO: 获取用户详细信息和位置
			Nickname: fmt.Sprintf("好友%d", f.FriendID),
			Avatar:   "",
		}
	}

	return resp, nil
}

// RemoveFriend 删除好友
func (s *FriendService) RemoveFriend(ctx context.Context, userID, friendID int64) error {
	return s.repo.RemoveFriend(ctx, userID, friendID)
}

// SearchUsers 搜索用户
func (s *FriendService) SearchUsers(ctx context.Context, userID int64, keyword string) ([]*dto.UserSearchResp, error) {
	// TODO: 实现用户搜索逻辑
	// 这里需要结合 user repository 和 model.User 进行搜索
	return []*dto.UserSearchResp{}, nil
}

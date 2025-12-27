package friend

import (
	"context"

	"gorm.io/gorm"

	"app/adaptor/repo/model"
)

// IFriendRepository 好友仓储接口
type IFriendRepository interface {
	CreateFriendRequest(ctx context.Context, req *model.FriendRequest) error
	GetFriendRequest(ctx context.Context, requestID int64) (*model.FriendRequest, error)
	GetPendingRequests(ctx context.Context, userID int64) ([]*model.FriendRequest, error)
	AcceptFriendRequest(ctx context.Context, requestID int64) error
	RejectFriendRequest(ctx context.Context, requestID int64) error
	CreateFriend(ctx context.Context, friend *model.Friend) error
	GetFriends(ctx context.Context, userID int64, status string) ([]*model.Friend, error)
	GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, error)
	RemoveFriend(ctx context.Context, userID, friendID int64) error
	AreFriends(ctx context.Context, userID, friendID int64) (bool, error)
}

// FriendRepository 好友仓储实现
type FriendRepository struct {
	db *gorm.DB
}

// NewFriendRepository 创建好友仓储
func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

// CreateFriendRequest 创建好友请求
func (r *FriendRepository) CreateFriendRequest(ctx context.Context, req *model.FriendRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

// GetFriendRequest 获取好友请求
func (r *FriendRepository) GetFriendRequest(ctx context.Context, requestID int64) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := r.db.WithContext(ctx).First(&req, requestID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &req, err
}

// GetPendingRequests 获取待处理的好友请求
func (r *FriendRepository) GetPendingRequests(ctx context.Context, userID int64) ([]*model.FriendRequest, error) {
	var reqs []*model.FriendRequest
	err := r.db.WithContext(ctx).Where("to_user_id = ? AND status = ?", userID, model.FriendStatusPending).
		Order("created_at DESC").Find(&reqs).Error
	return reqs, err
}

// AcceptFriendRequest 接受好友请求
func (r *FriendRepository) AcceptFriendRequest(ctx context.Context, requestID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新请求状态
		if err := tx.Model(&model.FriendRequest{}).Where("id = ?", requestID).
			Update("status", model.FriendStatusAccepted).Error; err != nil {
			return err
		}
		// 获取请求信息
		var req model.FriendRequest
		if err := tx.First(&req, requestID).Error; err != nil {
			return err
		}
		// 创建双向好友关系
		friend1 := &model.Friend{UserID: req.FromUserID, FriendID: req.ToUserID, Status: model.FriendStatusAccepted}
		friend2 := &model.Friend{UserID: req.ToUserID, FriendID: req.FromUserID, Status: model.FriendStatusAccepted}
		if err := tx.Create(friend1).Error; err != nil {
			return err
		}
		return tx.Create(friend2).Error
	})
}

// RejectFriendRequest 拒绝好友请求
func (r *FriendRepository) RejectFriendRequest(ctx context.Context, requestID int64) error {
	return r.db.WithContext(ctx).Model(&model.FriendRequest{}).Where("id = ?", requestID).
		Update("status", model.FriendStatusRejected).Error
}

// CreateFriend 创建好友关系
func (r *FriendRepository) CreateFriend(ctx context.Context, friend *model.Friend) error {
	return r.db.WithContext(ctx).Create(friend).Error
}

// GetFriends 获取好友列表
func (r *FriendRepository) GetFriends(ctx context.Context, userID int64, status string) ([]*model.Friend, error) {
	var friends []*model.Friend
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&friends).Error
	return friends, err
}

// GetFriend 获取特定好友关系
func (r *FriendRepository) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, error) {
	var friend model.Friend
	err := r.db.WithContext(ctx).Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &friend, err
}

// RemoveFriend 删除好友关系
func (r *FriendRepository) RemoveFriend(ctx context.Context, userID, friendID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除双向好友关系
		if err := tx.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ? AND friend_id = ?", friendID, userID).Delete(&model.Friend{}).Error
	})
}

// AreFriends 检查是否是好友
func (r *FriendRepository) AreFriends(ctx context.Context, userID, friendID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, model.FriendStatusAccepted).
		Count(&count).Error
	return count > 0, err
}

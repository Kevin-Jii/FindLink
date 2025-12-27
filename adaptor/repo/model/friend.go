package model

import (
	"time"
)

// FriendStatus 好友关系状态
type FriendStatus string

const (
	FriendStatusPending   FriendStatus = "pending"
	FriendStatusAccepted  FriendStatus = "accepted"
	FriendStatusRejected  FriendStatus = "rejected"
)

// SharingStatus 位置共享状态
type SharingStatus string

const (
	SharingStatusSharing SharingStatus = "sharing"
	SharingStatusPaused  SharingStatus = "paused"
	SharingStatusHidden  SharingStatus = "hidden"
)

// Friend 好友关系模型
type Friend struct {
	ID            int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64         `gorm:"not null;index" json:"user_id"`
	FriendID      int64         `gorm:"not null;index" json:"friend_id"`
	Status        FriendStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	SharingStatus SharingStatus `gorm:"type:varchar(20);default:'sharing'" json:"sharing_status"`
	CreatedAt     time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (*Friend) TableName() string {
	return "friends"
}

// FriendRequest 好友请求模型
type FriendRequest struct {
	ID          int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserID  int64        `gorm:"not null;index" json:"from_user_id"`
	ToUserID    int64        `gorm:"not null;index" json:"to_user_id"`
	Message     string       `gorm:"type:varchar(255)" json:"message"`
	Status      FriendStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

func (*FriendRequest) TableName() string {
	return "friend_requests"
}

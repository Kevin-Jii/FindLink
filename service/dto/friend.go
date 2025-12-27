package dto

import "time"

// FriendRequestReq 好友请求请求
type FriendRequestReq struct {
	ToUserID int64  `json:"to_user_id" binding:"required"`
	Message  string `json:"message"`
}

// FriendRequestActionReq 好友请求操作请求
type FriendRequestActionReq struct {
	RequestID int64 `json:"request_id" binding:"required"`
}

// FriendResp 好友响应
type FriendResp struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	FriendID      int64     `json:"friend_id"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Status        string    `json:"status"`         // pending, accepted, rejected
	SharingStatus string    `json:"sharing_status"` // sharing, paused, hidden
	LastActive    time.Time `json:"last_active"`
	Location      *LocationResp `json:"location,omitempty"`
}

// FriendRequestResp 好友请求响应
type FriendRequestResp struct {
	ID         int64     `json:"id"`
	FromUserID int64     `json:"from_user_id"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserSearchResp 用户搜索响应
type UserSearchResp struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
}

// FriendListReq 好友列表查询请求
type FriendListReq struct {
	Status string `json:"status"` // pending, accepted, rejected
}

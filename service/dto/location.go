package dto

import (
	"time"
)

// LocationMode 位置上报模式
type LocationMode string

const (
	LocationModeForeground        LocationMode = "foreground"
	LocationModeBackground        LocationMode = "background"
	LocationModeSignificantChange LocationMode = "significant_change"
)

// LocationReportReq 位置上报请求
type LocationReportReq struct {
	Longitude     float64      `json:"longitude" binding:"required"`
	Latitude      float64      `json:"latitude" binding:"required"`
	Accuracy      float64      `json:"accuracy"`
	Altitude      float64      `json:"altitude"`
	Speed         float64      `json:"speed"`
	Bearing       float64      `json:"bearing"`
	BatteryLevel  int          `json:"battery_level"`
	LocationMode  LocationMode `json:"location_mode"`
	IsLowAccuracy bool         `json:"is_low_accuracy"`
	DeviceID      string       `json:"device_id"` // 可选，设备上报时使用
}

// BatchLocationReportReq 批量位置上报请求
type BatchLocationReportReq struct {
	Locations []LocationReportReq `json:"locations" binding:"required,dive"`
}

// LocationResp 位置响应
type LocationResp struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id,omitempty"`
	DeviceID     string      `json:"device_id,omitempty"`
	Longitude    float64     `json:"longitude"`
	Latitude     float64     `json:"latitude"`
	Accuracy     float64     `json:"accuracy"`
	Altitude     float64     `json:"altitude"`
	Speed        float64     `json:"speed"`
	Bearing      float64     `json:"bearing"`
	BatteryLevel int         `json:"battery_level"`
	LocationMode LocationMode `json:"location_mode,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

// LocationHistoryReq 位置历史查询请求
type LocationHistoryReq struct {
	UserID   int64     `json:"user_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}

// NearbyFriendResp 附近好友响应
type NearbyFriendResp struct {
	UserID       int64       `json:"user_id"`
	Nickname     string      `json:"nickname"`
	Avatar       string      `json:"avatar"`
	Longitude    float64     `json:"longitude"`
	Latitude     float64     `json:"latitude"`
	Distance     float64     `json:"distance"` // 距离（米）
	BatteryLevel int         `json:"battery_level"`
	LastActive   time.Time   `json:"last_active"`
}

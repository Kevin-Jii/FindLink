package dto

import "time"

// DeviceBindReq 设备绑定请求
type DeviceBindReq struct {
	DeviceID string `json:"device_id" binding:"required"`
	Name     string `json:"name"`
	Type     string `json:"type"` // phone, tablet, watch, tracker, other
}

// DeviceResp 设备响应
type DeviceResp struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	BatteryLevel        int    `json:"battery_level"`
	ConnectionStatus    string `json:"connection_status"`
	NotificationEnabled bool   `json:"notification_enabled"`
	Location            *LocationResp `json:"location,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

// DeviceSettingsReq 设备设置请求
type DeviceSettingsReq struct {
	Name                string `json:"name"`
	NotificationEnabled *bool  `json:"notification_enabled"`
}

// DeviceStatusReq 设备状态更新请求
type DeviceStatusReq struct {
	BatteryLevel     int    `json:"battery_level"`
	ConnectionStatus string `json:"connection_status"` // online, offline, unknown
}

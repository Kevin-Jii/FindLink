package dto

import "time"

// UserSettingsReq 用户设置请求
type UserSettingsReq struct {
	ShareLocation *bool  `json:"share_location"`
	GhostMode     *bool  `json:"ghost_mode"`
	SmartAlerts   *bool  `json:"smart_alerts"`
	SOSAlerts     *bool  `json:"sos_alerts"`
	MapStyle      string `json:"map_style"`
	DistanceUnit  string `json:"distance_unit"`
}

// UserSettingsResp 用户设置响应
type UserSettingsResp struct {
	UserID        int64  `json:"user_id"`
	ShareLocation bool   `json:"share_location"`
	GhostMode     bool   `json:"ghost_mode"`
	SmartAlerts   bool   `json:"smart_alerts"`
	SOSAlerts     bool   `json:"sos_alerts"`
	MapStyle      string `json:"map_style"`
	DistanceUnit  string `json:"distance_unit"`
	UpdatedAt     time.Time `json:"updated_at"`
}

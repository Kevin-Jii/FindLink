package model

import (
	"time"
)

// MapStyle 地图样式
type MapStyle string

const (
	MapStyleDark  MapStyle = "dark"
	MapStyleLight MapStyle = "light"
	MapStyleStandard MapStyle = "standard"
)

// DistanceUnit 距离单位
type DistanceUnit string

const (
	DistanceUnitKm  DistanceUnit = "km"
	DistanceUnitMi  DistanceUnit = "mi"
)

// UserSettings 用户设置模型
type UserSettings struct {
	UserID        int64        `gorm:"primaryKey" json:"user_id"`
	ShareLocation bool         `gorm:"default:true" json:"share_location"`
	GhostMode     bool         `gorm:"default:false" json:"ghost_mode"`
	SmartAlerts   bool         `gorm:"default:true" json:"smart_alerts"`
	SOSAlerts     bool         `gorm:"default:true" json:"sos_alerts"`
	MapStyle      MapStyle     `gorm:"type:varchar(20);default:'dark'" json:"map_style"`
	DistanceUnit  DistanceUnit `gorm:"type:varchar(10);default:'km'" json:"distance_unit"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

func (*UserSettings) TableName() string {
	return "user_settings"
}

package model

import (
	"time"
)

// Device 设备模型
type Device struct {
	ID                  string                  `gorm:"primaryKey;type:varchar(64)" json:"id"`
	UserID              int64                   `gorm:"not null;index" json:"user_id"`
	Name                string                  `gorm:"type:varchar(100)" json:"name"`
	Type                string                  `gorm:"type:varchar(50);default:'unknown'" json:"type"`
	BatteryLevel        int                     `gorm:"type:int;default:0" json:"battery_level"`
	ConnectionStatus    DeviceConnectionStatus  `gorm:"type:varchar(20);default:'unknown'" json:"connection_status"`
	NotificationEnabled bool                    `gorm:"default:true" json:"notification_enabled"`
	CreatedAt           time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
}

func (*Device) TableName() string {
	return "devices"
}

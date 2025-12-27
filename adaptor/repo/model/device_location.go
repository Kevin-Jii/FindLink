package model

import (
	"time"
)

// DeviceConnectionStatus 设备连接状态
type DeviceConnectionStatus string

const (
	DeviceConnectionOnline  DeviceConnectionStatus = "online"
	DeviceConnectionOffline DeviceConnectionStatus = "offline"
	DeviceConnectionUnknown DeviceConnectionStatus = "unknown"
)

// DeviceLocation 设备位置模型
type DeviceLocation struct {
	ID            int64                   `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID      string                  `gorm:"not null;index" json:"device_id"`
	Location      string                  `gorm:"type:point;not null" json:"-"`
	Longitude     float64                 `gorm:"-" json:"longitude"`
	Latitude      float64                 `gorm:"-" json:"latitude"`
	Accuracy      float64                 `gorm:"type:float" json:"accuracy"`
	BatteryLevel  int                     `gorm:"type:int" json:"battery_level"`
	ConnectionStatus DeviceConnectionStatus `gorm:"type:varchar(20);default:'unknown'" json:"connection_status"`
	CreatedAt     time.Time               `gorm:"autoCreateTime" json:"created_at"`
}

func (*DeviceLocation) TableName() string {
	return "device_locations"
}

// ScanLocation 解析设备位置坐标
func (d *DeviceLocation) ScanLocation() error {
	if d.Location == "" {
		return nil
	}
	var lon, lat float64
	ok, err := scanPoint(d.Location, &lon, &lat)
	if !ok || err != nil {
		return err
	}
	d.Longitude = lon
	d.Latitude = lat
	return nil
}

// SetLocation 设置设备位置
func (d *DeviceLocation) SetLocation(lon, lat float64) {
	d.Location = "POINT(" + formatFloat(lon) + " " + formatFloat(lat) + ")"
	d.Longitude = lon
	d.Latitude = lat
}

package model

import (
	"time"
)

// Geofence 地理围栏模型
type Geofence struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"not null;index" json:"user_id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Center        string    `gorm:"type:point;not null" json:"-"`
	CenterLon     float64   `gorm:"-" json:"center_lon"`
	CenterLat     float64   `gorm:"-" json:"center_lat"`
	RadiusMeters  float64   `gorm:"type:float;not null" json:"radius_meters"`
	NotifyOnEnter bool      `gorm:"default:true" json:"notify_on_enter"`
	NotifyOnExit  bool      `gorm:"default:true" json:"notify_on_exit"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (*Geofence) TableName() string {
	return "geofences"
}

// ScanCenter 解析围栏中心坐标
func (g *Geofence) ScanCenter() error {
	if g.Center == "" {
		return nil
	}
	var lon, lat float64
	ok, err := scanPoint(g.Center, &lon, &lat)
	if !ok || err != nil {
		return err
	}
	g.CenterLon = lon
	g.CenterLat = lat
	return nil
}

// SetCenter 设置围栏中心
func (g *Geofence) SetCenter(lon, lat float64) {
	g.Center = "POINT(" + formatFloat(lon) + " " + formatFloat(lat) + ")"
	g.CenterLon = lon
	g.CenterLat = lat
}

// GeofenceEvent 地理围栏事件模型
type GeofenceEvent struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	GeofenceID  int64     `gorm:"not null;index" json:"geofence_id"`
	EntityType  string    `gorm:"type:varchar(20);not null" json:"entity_type"` // "user" or "device"
	EntityID    int64     `gorm:"not null" json:"entity_id"`
	EventType   string    `gorm:"type:varchar(20);not null" json:"event_type"` // "enter" or "exit"
	Location    string    `gorm:"type:point" json:"-"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (*GeofenceEvent) TableName() string {
	return "geofence_events"
}

package dto

import "time"

// GeofenceCreateReq 创建地理围栏请求
type GeofenceCreateReq struct {
	Name          string  `json:"name" binding:"required"`
	CenterLon     float64 `json:"center_lon" binding:"required"`
	CenterLat     float64 `json:"center_lat" binding:"required"`
	RadiusMeters  float64 `json:"radius_meters" binding:"required,gt=0"`
	NotifyOnEnter bool    `json:"notify_on_enter"`
	NotifyOnExit  bool    `json:"notify_on_exit"`
}

// GeofenceUpdateReq 更新地理围栏请求
type GeofenceUpdateReq struct {
	Name          string  `json:"name"`
	RadiusMeters  float64 `json:"radius_meters"`
	NotifyOnEnter *bool   `json:"notify_on_enter"`
	NotifyOnExit  *bool   `json:"notify_on_exit"`
	IsActive      *bool   `json:"is_active"`
}

// GeofenceResp 地理围栏响应
type GeofenceResp struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	CenterLon     float64   `json:"center_lon"`
	CenterLat     float64   `json:"center_lat"`
	RadiusMeters  float64   `json:"radius_meters"`
	NotifyOnEnter bool      `json:"notify_on_enter"`
	NotifyOnExit  bool      `json:"notify_on_exit"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GeofenceEventResp 地理围栏事件响应
type GeofenceEventResp struct {
	ID          int64     `json:"id"`
	GeofenceID  int64     `json:"geofence_id"`
	GeofenceName string  `json:"geofence_name"`
	EntityType  string    `json:"entity_type"`
	EntityID    int64     `json:"entity_id"`
	EventType   string    `json:"event_type"` // enter, exit
	CreatedAt   time.Time `json:"created_at"`
}

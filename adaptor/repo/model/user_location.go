package model

import (
	"time"
)

// LocationMode 位置上报模式
type LocationMode string

const (
	LocationModeForeground          LocationMode = "foreground"
	LocationModeBackground          LocationMode = "background"
	LocationModeSignificantChange   LocationMode = "significant_change"
)

// UserLocation 用户位置模型
type UserLocation struct {
	ID           int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64        `gorm:"not null;index" json:"user_id"`
	Location     string       `gorm:"type:point;not null" json:"-"` // MySQL POINT
	Longitude    float64      `gorm:"-" json:"longitude"`                           // 从 PostGIS 解析
	Latitude     float64      `gorm:"-" json:"latitude"`                            // 从 PostGIS 解析
	Accuracy     float64      `gorm:"type:float" json:"accuracy"`
	Altitude     float64      `gorm:"type:float" json:"altitude"`
	Speed        float64      `gorm:"type:float" json:"speed"`
	Bearing      float64      `gorm:"type:float" json:"bearing"`
	BatteryLevel int          `gorm:"type:int" json:"battery_level"`
	LocationMode LocationMode `gorm:"type:varchar(20);default:'foreground'" json:"location_mode"`
	IsLowAccuracy bool        `gorm:"default:false" json:"is_low_accuracy"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
}

func (*UserLocation) TableName() string {
	return "user_locations"
}

// ScanLocation 从 PostGIS geometry 解析坐标
func (u *UserLocation) ScanLocation() error {
	if u.Location == "" {
		return nil
	}
	// 解析 PostGIS ST_X 和 ST_Y
	// 格式: "POINT(longitude latitude)"
	var lon, lat float64
	_, err := scanPoint(u.Location, &lon, &lat)
	if err != nil {
		return err
	}
	u.Longitude = lon
	u.Latitude = lat
	return nil
}

// SetLocation 设置 PostGIS POINT
func (u *UserLocation) SetLocation(lon, lat float64) {
	u.Location = u.makePoint(lon, lat)
	u.Longitude = lon
	u.Latitude = lat
}

func (u *UserLocation) makePoint(lon, lat float64) string {
	return "POINT(" + formatFloat(lon) + " " + formatFloat(lat) + ")"
}

func scanPoint(wkt string, lon, lat *float64) (bool, error) {
	// 简单解析 "POINT(lon lat)" 格式
	if len(wkt) < 8 || wkt[0:5] != "POINT" {
		return false, nil
	}
	var l1, l2 float64
	_, err := scanFloat(wkt[6:], &l1)
	if err != nil {
		return false, err
	}
	_, err = scanFloat(wkt[6:], &l2)
	if err != nil {
		return false, err
	}
	// 找到空格分隔的两个值
	idx := findSpace(wkt)
	if idx == -1 {
		return false, nil
	}
	_, err = scanFloat(wkt[6:idx], lon)
	if err != nil {
		return false, err
	}
	_, err = scanFloat(wkt[idx+1:len(wkt)-1], lat)
	if err != nil {
		return false, err
	}
	return true, nil
}

func scanFloat(s string, f *float64) (int, error) {
	var val float64
	neg := false
	start := 0
	if len(s) > 0 && s[0] == '-' {
		neg = true
		start = 1
	}
	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			if s[i] == '.' {
				continue
			}
			break
		}
		val = val*10 + float64(s[i]-'0')
	}
	if neg {
		val = -val
	}
	*f = val
	return len(s), nil
}

func findSpace(s string) int {
	for i := 6; i < len(s); i++ {
		if s[i] == ' ' {
			return i
		}
	}
	return -1
}

func formatFloat(f float64) string {
	return time.Duration(f * 1e9).String()[:len(time.Duration(f*1e9).String())-1]
}

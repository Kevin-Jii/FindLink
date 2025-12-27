package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"app/service/dto"
)

const (
	// Location cache keys
	userLocationKey    = "loc:user:%d"
	deviceLocationKey  = "loc:device:%s"
	friendsKey         = "friends:%d"
	userSettingsKey    = "settings:%d"

	// GEO keys
	geoUsersKey   = "geo:users"
	geoDevicesKey = "geo:devices"

	// Cache TTLs
	locationTTL = 5 * time.Minute
	friendsTTL  = 1 * time.Hour
	settingsTTL = 1 * time.Hour
)

// LocationCache 位置缓存服务
type LocationCache struct {
	client *redis.Client
}

// NewLocationCache 创建位置缓存服务
func NewLocationCache(client *redis.Client) *LocationCache {
	return &LocationCache{client: client}
}

// SetUserLocation 设置用户位置
func (c *LocationCache) SetUserLocation(userID int64, loc *dto.LocationResp) error {
	data, err := json.Marshal(loc)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(userLocationKey, userID)

	pipe := c.client.Pipeline()
	pipe.Set(key, data, locationTTL)
	pipe.GeoAdd(geoUsersKey, &redis.GeoLocation{
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
		Name:      fmt.Sprintf("%d", userID),
	})
	_, err = pipe.Exec()
	return err
}

// GetUserLocation 获取用户位置
func (c *LocationCache) GetUserLocation(userID int64) (*dto.LocationResp, error) {
	key := fmt.Sprintf(userLocationKey, userID)
	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var loc dto.LocationResp
	if err := json.Unmarshal(data, &loc); err != nil {
		return nil, err
	}
	return &loc, nil
}

// SetDeviceLocation 设置设备位置
func (c *LocationCache) SetDeviceLocation(deviceID string, loc *dto.LocationResp) error {
	data, err := json.Marshal(loc)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(deviceLocationKey, deviceID)

	pipe := c.client.Pipeline()
	pipe.Set(key, data, locationTTL)
	pipe.GeoAdd(geoDevicesKey, &redis.GeoLocation{
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
		Name:      deviceID,
	})
	_, err = pipe.Exec()
	return err
}

// GetDeviceLocation 获取设备位置
func (c *LocationCache) GetDeviceLocation(deviceID string) (*dto.LocationResp, error) {
	key := fmt.Sprintf(deviceLocationKey, deviceID)
	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var loc dto.LocationResp
	if err := json.Unmarshal(data, &loc); err != nil {
		return nil, err
	}
	return &loc, nil
}

// DeleteUserLocation 删除用户位置缓存
func (c *LocationCache) DeleteUserLocation(userID int64) error {
	key := fmt.Sprintf(userLocationKey, userID)
	// 使用 ZREM 删除 GEO set 中的成员
	c.client.ZRem(geoUsersKey, fmt.Sprintf("%d", userID))
	return c.client.Del(key).Err()
}

// DeleteDeviceLocation 删除设备位置缓存
func (c *LocationCache) DeleteDeviceLocation(deviceID string) error {
	key := fmt.Sprintf(deviceLocationKey, deviceID)
	c.client.ZRem(geoDevicesKey, deviceID)
	return c.client.Del(key).Err()
}

// GeoRadius GEO半径查询
func (c *LocationCache) GeoRadius(key string, lon, lat, radiusMeters float64) ([]string, error) {
	locations, err := c.client.GeoRadius(key, lon, lat, &redis.GeoRadiusQuery{
		Radius:      radiusMeters,
		Unit:        "m",
		WithCoord:   false,
		WithDist:    false,
		WithGeoHash: false,
		Count:       0,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(locations))
	for i, loc := range locations {
		names[i] = loc.Name
	}
	return names, nil
}

// GeoRadiusUsers 查询附近用户
func (c *LocationCache) GeoRadiusUsers(lon, lat, radiusMeters float64) ([]string, error) {
	return c.GeoRadius(geoUsersKey, lon, lat, radiusMeters)
}

// GeoRadiusDevices 查询附近设备
func (c *LocationCache) GeoRadiusDevices(lon, lat, radiusMeters float64) ([]string, error) {
	return c.GeoRadius(geoDevicesKey, lon, lat, radiusMeters)
}

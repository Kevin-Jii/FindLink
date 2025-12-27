package adaptor

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"app/config"
	"app/adaptor/repo/device"
	"app/adaptor/repo/friend"
	"app/adaptor/repo/geofence"
	"app/adaptor/repo/location"
	redisCache "app/adaptor/redis"
)

type IAdaptor interface {
	GetConfig() *config.Config
	GetDB() *gorm.DB
	GetRedis() *redis.Client

	// 位置缓存
	NewLocationCache() *redisCache.LocationCache

	// 仓储
	NewLocationRepository() *location.LocationRepository
	NewFriendRepository() *friend.FriendRepository
	NewDeviceRepository() *device.DeviceRepository
	NewGeofenceRepository() *geofence.GeofenceRepository
}

type Adaptor struct {
	conf  *config.Config
	db    *gorm.DB
	redis *redis.Client
}

func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		conf:  conf,
		db:    db,
		redis: redis,
	}
}

func (a *Adaptor) GetConfig() *config.Config {
	return a.conf
}

func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}

func (a *Adaptor) GetRedis() *redis.Client {
	return a.redis
}

// 位置缓存
func (a *Adaptor) NewLocationCache() *redisCache.LocationCache {
	return redisCache.NewLocationCache(a.redis)
}

// 仓储
func (a *Adaptor) NewLocationRepository() *location.LocationRepository {
	return location.NewLocationRepository(a.db)
}

func (a *Adaptor) NewFriendRepository() *friend.FriendRepository {
	return friend.NewFriendRepository(a.db)
}

func (a *Adaptor) NewDeviceRepository() *device.DeviceRepository {
	return device.NewDeviceRepository(a.db)
}

func (a *Adaptor) NewGeofenceRepository() *geofence.GeofenceRepository {
	return geofence.NewGeofenceRepository(a.db)
}

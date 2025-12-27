package location

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"app/adaptor/repo/model"
)

// ILocationRepository 位置仓储接口
type ILocationRepository interface {
	CreateUserLocation(ctx context.Context, loc *model.UserLocation) error
	BatchCreateUserLocations(ctx context.Context, locs []*model.UserLocation) error
	GetLatestUserLocation(ctx context.Context, userID int64) (*model.UserLocation, error)
	GetUserLocationHistory(ctx context.Context, userID int64, startTime, endTime time.Time, limit, offset int) ([]*model.UserLocation, error)
	CreateDeviceLocation(ctx context.Context, loc *model.DeviceLocation) error
	GetLatestDeviceLocation(ctx context.Context, deviceID string) (*model.DeviceLocation, error)
	DeleteUserLocations(ctx context.Context, userID int64) error
}

// LocationRepository 位置仓储实现
type LocationRepository struct {
	db *gorm.DB
}

// NewLocationRepository 创建位置仓储
func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// CreateUserLocation 创建用户位置
func (r *LocationRepository) CreateUserLocation(ctx context.Context, loc *model.UserLocation) error {
	return r.db.WithContext(ctx).Create(loc).Error
}

// BatchCreateUserLocations 批量创建用户位置
func (r *LocationRepository) BatchCreateUserLocations(ctx context.Context, locs []*model.UserLocation) error {
	if len(locs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoNothing: true}).CreateInBatches(locs, 100).Error
}

// GetLatestUserLocation 获取用户最新位置
func (r *LocationRepository) GetLatestUserLocation(ctx context.Context, userID int64) (*model.UserLocation, error) {
	var loc model.UserLocation
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").First(&loc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	loc.ScanLocation()
	return &loc, nil
}

// GetUserLocationHistory 获取用户位置历史
func (r *LocationRepository) GetUserLocationHistory(ctx context.Context, userID int64, startTime, endTime time.Time, limit, offset int) ([]*model.UserLocation, error) {
	var locs []*model.UserLocation
	query := r.db.WithContext(ctx).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startTime, endTime)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Order("created_at DESC").Find(&locs).Error
	if err != nil {
		return nil, err
	}
	// 解析所有位置的坐标
	for _, loc := range locs {
		loc.ScanLocation()
	}
	return locs, nil
}

// CreateDeviceLocation 创建设备位置
func (r *LocationRepository) CreateDeviceLocation(ctx context.Context, loc *model.DeviceLocation) error {
	return r.db.WithContext(ctx).Create(loc).Error
}

// GetLatestDeviceLocation 获取设备最新位置
func (r *LocationRepository) GetLatestDeviceLocation(ctx context.Context, deviceID string) (*model.DeviceLocation, error) {
	var loc model.DeviceLocation
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("created_at DESC").First(&loc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	loc.ScanLocation()
	return &loc, nil
}

// DeleteUserLocations 删除用户位置记录
func (r *LocationRepository) DeleteUserLocations(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserLocation{}).Error
}

// GeoSearchNearby 使用 PostGIS 搜索附近位置
func (r *LocationRepository) GeoSearchNearby(ctx context.Context, lon, lat, radiusMeters float64, userIDs []int64) ([]*model.UserLocation, error) {
	var locs []*model.UserLocation
	query := fmt.Sprintf(`
		SELECT * FROM user_locations
		WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography, %f)
		ORDER BY created_at DESC
	`, lon, lat, radiusMeters)

	if len(userIDs) > 0 {
		query += fmt.Sprintf(" AND user_id IN (%v)", userIDs)
	}

	err := r.db.WithContext(ctx).Raw(query).Scan(&locs).Error
	if err != nil {
		return nil, err
	}
	return locs, nil
}

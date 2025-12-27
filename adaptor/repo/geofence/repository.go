package geofence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"app/adaptor/repo/model"
)

// IGeofenceRepository 地理围栏仓储接口
type IGeofenceRepository interface {
	Create(ctx context.Context, geofence *model.Geofence) error
	Get(ctx context.Context, geofenceID int64) (*model.Geofence, error)
	ListByUser(ctx context.Context, userID int64) ([]*model.Geofence, error)
	Update(ctx context.Context, geofence *model.Geofence) error
	Delete(ctx context.Context, geofenceID int64) error
	GetActiveGeofences(ctx context.Context, userID int64) ([]*model.Geofence, error)
	CheckPointInGeofences(ctx context.Context, lon, lat float64, entityIDs []int64) ([]*model.Geofence, error)
	CreateEvent(ctx context.Context, event *model.GeofenceEvent) error
}

// GeofenceRepository 地理围栏仓储实现
type GeofenceRepository struct {
	db *gorm.DB
}

// NewGeofenceRepository 创建地理围栏仓储
func NewGeofenceRepository(db *gorm.DB) *GeofenceRepository {
	return &GeofenceRepository{db: db}
}

// Create 创建地理围栏
func (r *GeofenceRepository) Create(ctx context.Context, geofence *model.Geofence) error {
	return r.db.WithContext(ctx).Create(geofence).Error
}

// Get 获取地理围栏
func (r *GeofenceRepository) Get(ctx context.Context, geofenceID int64) (*model.Geofence, error) {
	var geofence model.Geofence
	err := r.db.WithContext(ctx).First(&geofence, geofenceID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	geofence.ScanCenter()
	return &geofence, nil
}

// ListByUser 获取用户的所有地理围栏
func (r *GeofenceRepository) ListByUser(ctx context.Context, userID int64) ([]*model.Geofence, error) {
	var geofences []*model.Geofence
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&geofences).Error
	if err != nil {
		return nil, err
	}
	for _, g := range geofences {
		g.ScanCenter()
	}
	return geofences, nil
}

// Update 更新地理围栏
func (r *GeofenceRepository) Update(ctx context.Context, geofence *model.Geofence) error {
	return r.db.WithContext(ctx).Save(geofence).Error
}

// Delete 删除地理围栏
func (r *GeofenceRepository) Delete(ctx context.Context, geofenceID int64) error {
	return r.db.WithContext(ctx).Delete(&model.Geofence{}, geofenceID).Error
}

// GetActiveGeofences 获取用户活跃的地理围栏
func (r *GeofenceRepository) GetActiveGeofences(ctx context.Context, userID int64) ([]*model.Geofence, error) {
	var geofences []*model.Geofence
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&geofences).Error
	if err != nil {
		return nil, err
	}
	for _, g := range geofences {
		g.ScanCenter()
	}
	return geofences, nil
}

// CheckPointInGeofences 检查点是否在围栏内
func (r *GeofenceRepository) CheckPointInGeofences(ctx context.Context, lon, lat float64, entityIDs []int64) ([]*model.Geofence, error) {
	var geofences []*model.Geofence
	query := fmt.Sprintf(`
		SELECT * FROM geofences
		WHERE is_active = true
		AND ST_DWithin(center, ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography, radius_meters)
	`, lon, lat)

	if len(entityIDs) > 0 {
		query += fmt.Sprintf(" AND user_id IN (%v)", entityIDs)
	}

	err := r.db.WithContext(ctx).Raw(query).Scan(&geofences).Error
	if err != nil {
		return nil, err
	}
	return geofences, nil
}

// CreateEvent 创建围栏事件
func (r *GeofenceRepository) CreateEvent(ctx context.Context, event *model.GeofenceEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

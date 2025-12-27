package device

import (
	"context"

	"gorm.io/gorm"

	"app/adaptor/repo/model"
)

// IDeviceRepository 设备仓储接口
type IDeviceRepository interface {
	Create(ctx context.Context, device *model.Device) error
	Get(ctx context.Context, deviceID string) (*model.Device, error)
	GetByUser(ctx context.Context, userID int64) ([]*model.Device, error)
	Update(ctx context.Context, device *model.Device) error
	UpdateStatus(ctx context.Context, deviceID string, batteryLevel int, connectionStatus model.DeviceConnectionStatus) error
	Delete(ctx context.Context, deviceID string) error
	IsBound(ctx context.Context, deviceID string) (bool, error)
	GetUserByDevice(ctx context.Context, deviceID string) (int64, error)
}

// DeviceRepository 设备仓储实现
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository 创建设备仓储
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create 绑定设备
func (r *DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// Get 获取设备
func (r *DeviceRepository) Get(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &device, err
}

// GetByUser 获取用户的设备列表
func (r *DeviceRepository) GetByUser(ctx context.Context, userID int64) ([]*model.Device, error) {
	var devices []*model.Device
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&devices).Error
	return devices, err
}

// Update 更新设备
func (r *DeviceRepository) Update(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// UpdateStatus 更新设备状态
func (r *DeviceRepository) UpdateStatus(ctx context.Context, deviceID string, batteryLevel int, connectionStatus model.DeviceConnectionStatus) error {
	return r.db.WithContext(ctx).Model(&model.Device{}).
		Where("id = ?", deviceID).
		Updates(map[string]interface{}{
			"battery_level":      batteryLevel,
			"connection_status":  connectionStatus,
		}).Error
}

// Delete 解绑设备
func (r *DeviceRepository) Delete(ctx context.Context, deviceID string) error {
	return r.db.WithContext(ctx).Where("id = ?", deviceID).Delete(&model.Device{}).Error
}

// IsBound 检查设备是否已绑定
func (r *DeviceRepository) IsBound(ctx context.Context, deviceID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", deviceID).Count(&count).Error
	return count > 0, err
}

// GetUserByDevice 获取设备绑定的用户ID
func (r *DeviceRepository) GetUserByDevice(ctx context.Context, deviceID string) (int64, error) {
	var device model.Device
	err := r.db.WithContext(ctx).Select("user_id").First(&device, "id = ?", deviceID).Error
	if err != nil {
		return 0, err
	}
	return device.UserID, nil
}

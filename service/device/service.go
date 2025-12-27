package device

import (
	"context"

	"app/adaptor/repo/device"
	"app/adaptor/repo/model"
	"app/common"
	"app/service/dto"
)

// IDeviceService 设备服务接口
type IDeviceService interface {
	BindDevice(ctx context.Context, userID int64, req *dto.DeviceBindReq) error
	UnbindDevice(ctx context.Context, userID int64, deviceID string) error
	GetDeviceList(ctx context.Context, userID int64) ([]*dto.DeviceResp, error)
	UpdateDeviceSettings(ctx context.Context, userID int64, deviceID string, req *dto.DeviceSettingsReq) error
	UpdateDeviceStatus(ctx context.Context, deviceID string, req *dto.DeviceStatusReq) error
}

// DeviceService 设备服务实现
type DeviceService struct {
	repo *device.DeviceRepository
}

// NewDeviceService 创建设备服务
func NewDeviceService(repo *device.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// BindDevice 绑定设备
func (s *DeviceService) BindDevice(ctx context.Context, userID int64, req *dto.DeviceBindReq) error {
	// 检查设备是否已被绑定
	bound, err := s.repo.IsBound(ctx, req.DeviceID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if bound {
		return common.DeviceAlreadyBoundErr
	}

	deviceType := req.Type
	if deviceType == "" {
		deviceType = "other"
	}

	d := &model.Device{
		ID:                  req.DeviceID,
		UserID:              userID,
		Name:                req.Name,
		Type:                deviceType,
		ConnectionStatus:    model.DeviceConnectionUnknown,
		NotificationEnabled: true,
	}

	return s.repo.Create(ctx, d)
}

// UnbindDevice 解绑设备
func (s *DeviceService) UnbindDevice(ctx context.Context, userID int64, deviceID string) error {
	// 检查设备是否存在
	d, err := s.repo.Get(ctx, deviceID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if d == nil {
		return common.DeviceNotFoundErr
	}

	// 检查设备是否属于该用户
	if d.UserID != userID {
		return common.PermissionErr
	}

	return s.repo.Delete(ctx, deviceID)
}

// GetDeviceList 获取设备列表
func (s *DeviceService) GetDeviceList(ctx context.Context, userID int64) ([]*dto.DeviceResp, error) {
	devices, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}

	resp := make([]*dto.DeviceResp, len(devices))
	for i, d := range devices {
		resp[i] = &dto.DeviceResp{
			ID:                  d.ID,
			Name:                d.Name,
			Type:                d.Type,
			BatteryLevel:        d.BatteryLevel,
			ConnectionStatus:    string(d.ConnectionStatus),
			NotificationEnabled: d.NotificationEnabled,
			CreatedAt:           d.CreatedAt,
		}
	}

	return resp, nil
}

// UpdateDeviceSettings 更新设备设置
func (s *DeviceService) UpdateDeviceSettings(ctx context.Context, userID int64, deviceID string, req *dto.DeviceSettingsReq) error {
	d, err := s.repo.Get(ctx, deviceID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if d == nil {
		return common.DeviceNotFoundErr
	}
	if d.UserID != userID {
		return common.PermissionErr
	}

	if req.Name != "" {
		d.Name = req.Name
	}
	if req.NotificationEnabled != nil {
		d.NotificationEnabled = *req.NotificationEnabled
	}

	return s.repo.Update(ctx, d)
}

// UpdateDeviceStatus 更新设备状态
func (s *DeviceService) UpdateDeviceStatus(ctx context.Context, deviceID string, req *dto.DeviceStatusReq) error {
	status := model.DeviceConnectionStatus(req.ConnectionStatus)
	if status == "" {
		status = model.DeviceConnectionUnknown
	}

	return s.repo.UpdateStatus(ctx, deviceID, req.BatteryLevel, status)
}

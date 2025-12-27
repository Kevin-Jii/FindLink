package geofence

import (
	"context"

	"app/adaptor/repo/geofence"
	"app/adaptor/repo/model"
	"app/common"
	"app/service/dto"
)

// IGeofenceService 地理围栏服务接口
type IGeofenceService interface {
	Create(ctx context.Context, userID int64, req *dto.GeofenceCreateReq) error
	GetList(ctx context.Context, userID int64) ([]*dto.GeofenceResp, error)
	Update(ctx context.Context, userID int64, geofenceID int64, req *dto.GeofenceUpdateReq) error
	Delete(ctx context.Context, userID int64, geofenceID int64) error
	CheckGeofenceEvents(ctx context.Context, lon, lat float64, entityType string, entityID int64) error
}

// GeofenceService 地理围栏服务实现
type GeofenceService struct {
	repo *geofence.GeofenceRepository
}

// NewGeofenceService 创建地理围栏服务
func NewGeofenceService(repo *geofence.GeofenceRepository) *GeofenceService {
	return &GeofenceService{repo: repo}
}

// Create 创建地理围栏
func (s *GeofenceService) Create(ctx context.Context, userID int64, req *dto.GeofenceCreateReq) error {
	g := &model.Geofence{
		UserID:        userID,
		Name:          req.Name,
		RadiusMeters:  req.RadiusMeters,
		NotifyOnEnter: req.NotifyOnEnter,
		NotifyOnExit:  req.NotifyOnExit,
		IsActive:      true,
	}
	g.SetCenter(req.CenterLon, req.CenterLat)

	return s.repo.Create(ctx, g)
}

// GetList 获取地理围栏列表
func (s *GeofenceService) GetList(ctx context.Context, userID int64) ([]*dto.GeofenceResp, error) {
	geofences, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}

	resp := make([]*dto.GeofenceResp, len(geofences))
	for i, g := range geofences {
		resp[i] = &dto.GeofenceResp{
			ID:            g.ID,
			Name:          g.Name,
			CenterLon:     g.CenterLon,
			CenterLat:     g.CenterLat,
			RadiusMeters:  g.RadiusMeters,
			NotifyOnEnter: g.NotifyOnEnter,
			NotifyOnExit:  g.NotifyOnExit,
			IsActive:      g.IsActive,
			CreatedAt:     g.CreatedAt,
			UpdatedAt:     g.UpdatedAt,
		}
	}

	return resp, nil
}

// Update 更新地理围栏
func (s *GeofenceService) Update(ctx context.Context, userID int64, geofenceID int64, req *dto.GeofenceUpdateReq) error {
	g, err := s.repo.Get(ctx, geofenceID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if g == nil {
		return common.GeofenceNotFoundErr
	}
	if g.UserID != userID {
		return common.PermissionErr
	}

	if req.Name != "" {
		g.Name = req.Name
	}
	if req.RadiusMeters > 0 {
		g.RadiusMeters = req.RadiusMeters
	}
	if req.NotifyOnEnter != nil {
		g.NotifyOnEnter = *req.NotifyOnEnter
	}
	if req.NotifyOnExit != nil {
		g.NotifyOnExit = *req.NotifyOnExit
	}
	if req.IsActive != nil {
		g.IsActive = *req.IsActive
	}

	return s.repo.Update(ctx, g)
}

// Delete 删除地理围栏
func (s *GeofenceService) Delete(ctx context.Context, userID int64, geofenceID int64) error {
	g, err := s.repo.Get(ctx, geofenceID)
	if err != nil {
		return common.DatabaseErr.WithErr(err)
	}
	if g == nil {
		return common.GeofenceNotFoundErr
	}
	if g.UserID != userID {
		return common.PermissionErr
	}

	return s.repo.Delete(ctx, geofenceID)
}

// CheckGeofenceEvents 检查围栏事件
func (s *GeofenceService) CheckGeofenceEvents(ctx context.Context, lon, lat float64, entityType string, entityID int64) error {
	// TODO: 实现围栏事件检测逻辑
	// 需要记录上次位置，判断是进入还是离开
	return nil
}

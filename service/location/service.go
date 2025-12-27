package location

import (
	"context"
	"fmt"

	"app/adaptor/repo/location"
	"app/adaptor/repo/model"
	redisCache "app/adaptor/redis"
	"app/common"
	"app/service/dto"
)

// ILocationService 位置服务接口
type ILocationService interface {
	ReportLocation(ctx context.Context, userID int64, req *dto.LocationReportReq) error
	BatchReportLocation(ctx context.Context, userID int64, req *dto.BatchLocationReportReq) error
	GetUserLocation(ctx context.Context, userID int64, requesterID int64) (*dto.LocationResp, error)
	GetDeviceLocation(ctx context.Context, deviceID string, userID int64) (*dto.LocationResp, error)
	GetLocationHistory(ctx context.Context, req *dto.LocationHistoryReq) ([]*dto.LocationResp, error)
	GetNearbyFriends(ctx context.Context, userID int64, radiusMeters float64) ([]*dto.NearbyFriendResp, error)
}

// LocationService 位置服务实现
type LocationService struct {
	repo  *location.LocationRepository
	cache *redisCache.LocationCache
}

// NewLocationService 创建位置服务
func NewLocationService(repo *location.LocationRepository, cache *redisCache.LocationCache) *LocationService {
	return &LocationService{repo: repo, cache: cache}
}

// ReportLocation 上报位置
func (s *LocationService) ReportLocation(ctx context.Context, userID int64, req *dto.LocationReportReq) error {
	// 标记低精度位置
	if req.Accuracy > 0 && req.Accuracy < 100 {
		req.IsLowAccuracy = true
	}

	loc := &model.UserLocation{
		UserID:        userID,
		Accuracy:      req.Accuracy,
		Altitude:      req.Altitude,
		Speed:         req.Speed,
		Bearing:       req.Bearing,
		BatteryLevel:  req.BatteryLevel,
		LocationMode:  model.LocationMode(req.LocationMode),
		IsLowAccuracy: req.IsLowAccuracy,
	}
	loc.SetLocation(req.Longitude, req.Latitude)

	// 保存到数据库
	if err := s.repo.CreateUserLocation(ctx, loc); err != nil {
		return common.DatabaseErr.WithErr(err)
	}

	// 更新缓存
	resp := s.toLocationResp(loc)
	if err := s.cache.SetUserLocation(userID, resp); err != nil {
		// 缓存失败不影响业务
		fmt.Printf("cache user location failed: %v\n", err)
	}

	return nil
}

// BatchReportLocation 批量上报位置
func (s *LocationService) BatchReportLocation(ctx context.Context, userID int64, req *dto.BatchLocationReportReq) error {
	if len(req.Locations) == 0 {
		return nil
	}

	locs := make([]*model.UserLocation, 0, len(req.Locations))
	for _, locReq := range req.Locations {
		if locReq.Accuracy > 0 && locReq.Accuracy < 100 {
			locReq.IsLowAccuracy = true
		}

		loc := &model.UserLocation{
			UserID:        userID,
			Accuracy:      locReq.Accuracy,
			Altitude:      locReq.Altitude,
			Speed:         locReq.Speed,
			Bearing:       locReq.Bearing,
			BatteryLevel:  locReq.BatteryLevel,
			LocationMode:  model.LocationMode(locReq.LocationMode),
			IsLowAccuracy: locReq.IsLowAccuracy,
		}
		loc.SetLocation(locReq.Longitude, locReq.Latitude)
		locs = append(locs, loc)
	}

	if err := s.repo.BatchCreateUserLocations(ctx, locs); err != nil {
		return common.DatabaseErr.WithErr(err)
	}

	// 更新缓存为最新位置
	if len(locs) > 0 {
		resp := s.toLocationResp(locs[len(locs)-1])
		if err := s.cache.SetUserLocation(userID, resp); err != nil {
			fmt.Printf("cache user location failed: %v\n", err)
		}
	}

	return nil
}

// GetUserLocation 获取用户位置
func (s *LocationService) GetUserLocation(ctx context.Context, userID int64, requesterID int64) (*dto.LocationResp, error) {
	// 先从缓存获取
	resp, err := s.cache.GetUserLocation(userID)
	if err != nil {
		return nil, common.RedisErr.WithErr(err)
	}
	if resp != nil {
		return resp, nil
	}

	// 缓存未命中，从数据库获取
	loc, err := s.repo.GetLatestUserLocation(ctx, userID)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}
	if loc == nil {
		return nil, common.LocationNotFoundErr
	}

	resp = s.toLocationResp(loc)
	// 写入缓存
	if err := s.cache.SetUserLocation(userID, resp); err != nil {
		fmt.Printf("cache user location failed: %v\n", err)
	}

	return resp, nil
}

// GetDeviceLocation 获取设备位置
func (s *LocationService) GetDeviceLocation(ctx context.Context, deviceID string, userID int64) (*dto.LocationResp, error) {
	// 先从缓存获取
	resp, err := s.cache.GetDeviceLocation(deviceID)
	if err != nil {
		return nil, common.RedisErr.WithErr(err)
	}
	if resp != nil {
		return resp, nil
	}

	// 缓存未命中，从数据库获取
	loc, err := s.repo.GetLatestDeviceLocation(ctx, deviceID)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}
	if loc == nil {
		return nil, common.LocationNotFoundErr
	}

	resp = &dto.LocationResp{
		DeviceID:     deviceID,
		Longitude:    loc.Longitude,
		Latitude:     loc.Latitude,
		Accuracy:     loc.Accuracy,
		BatteryLevel: loc.BatteryLevel,
		CreatedAt:    loc.CreatedAt,
	}

	// 写入缓存
	if err := s.cache.SetDeviceLocation(deviceID, resp); err != nil {
		fmt.Printf("cache device location failed: %v\n", err)
	}

	return resp, nil
}

// GetLocationHistory 获取位置历史
func (s *LocationService) GetLocationHistory(ctx context.Context, req *dto.LocationHistoryReq) ([]*dto.LocationResp, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	locs, err := s.repo.GetUserLocationHistory(ctx, req.UserID, req.StartTime, req.EndTime, limit, req.Offset)
	if err != nil {
		return nil, common.DatabaseErr.WithErr(err)
	}

	resp := make([]*dto.LocationResp, len(locs))
	for i, loc := range locs {
		resp[i] = s.toLocationResp(loc)
	}

	return resp, nil
}

// GetNearbyFriends 获取附近好友
func (s *LocationService) GetNearbyFriends(ctx context.Context, userID int64, radiusMeters float64) ([]*dto.NearbyFriendResp, error) {
	// 先尝试从 Redis GEO 获取
	friendIDs, err := s.cache.GeoRadiusUsers(0, 0, radiusMeters)
	if err != nil {
		return nil, common.RedisErr.WithErr(err)
	}

	// 如果 Redis GEO 没有数据，使用数据库查询
	if len(friendIDs) == 0 {
		// 这里需要结合好友关系查询
		// 简化处理，返回空列表
		return []*dto.NearbyFriendResp{}, nil
	}

	// TODO: 根据好友ID获取详细信息和位置
	// 这里需要结合 friend repository 和 user service

	return []*dto.NearbyFriendResp{}, nil
}

// toLocationResp 转换为响应
func (s *LocationService) toLocationResp(loc *model.UserLocation) *dto.LocationResp {
	return &dto.LocationResp{
		ID:            loc.ID,
		UserID:        loc.UserID,
		Longitude:     loc.Longitude,
		Latitude:      loc.Latitude,
		Accuracy:      loc.Accuracy,
		Altitude:      loc.Altitude,
		Speed:         loc.Speed,
		Bearing:       loc.Bearing,
		BatteryLevel:  loc.BatteryLevel,
		LocationMode:  dto.LocationMode(loc.LocationMode),
		CreatedAt:     loc.CreatedAt,
	}
}

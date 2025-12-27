package customer

import (
	"app/adaptor"
	"app/service/device"
	"app/service/friend"
	"app/service/geofence"
	"app/service/location"
	"app/service/user"
	"app/service/websocket"
)

type Ctrl struct {
	Adaptor  adaptor.IAdaptor
	User     *user.Service
	Location *location.LocationService
	Friend   *friend.FriendService
	Device   *device.DeviceService
	Geofence *geofence.GeofenceService
	Hub      *websocket.Hub
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	// 初始化Redis缓存
	locationCache := adaptor.NewLocationCache()

	// 初始化仓储
	locationRepo := adaptor.NewLocationRepository()
	friendRepo := adaptor.NewFriendRepository()
	deviceRepo := adaptor.NewDeviceRepository()
	geofenceRepo := adaptor.NewGeofenceRepository()

	// 初始化服务
	locationSvc := location.NewLocationService(locationRepo, locationCache)
	friendSvc := friend.NewFriendService(friendRepo)
	deviceSvc := device.NewDeviceService(deviceRepo)
	geofenceSvc := geofence.NewGeofenceService(geofenceRepo)

	// 初始化WebSocket Hub
	hub := websocket.NewHub()

	// 启动Hub
	go hub.Run()

	return &Ctrl{
		Adaptor:  adaptor,
		User:     user.NewService(adaptor),
		Location: locationSvc,
		Friend:   friendSvc,
		Device:   deviceSvc,
		Geofence: geofenceSvc,
		Hub:      hub,
	}
}

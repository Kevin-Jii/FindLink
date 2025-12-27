package common

type Errno struct {
	Code   int
	Msg    string
	ErrMsg string
}

func (err Errno) Error() string {
	return err.Msg
}

func (err Errno) WithMsg(msg string) Errno {
	err.Msg = err.Msg + "," + msg
	return err
}

func (err Errno) WithErr(rawErr error) Errno {
	var msg string
	if rawErr != nil {
		msg = rawErr.Error()
	}
	err.ErrMsg = err.Msg + "," + msg
	return err
}

func (err Errno) IsOk() bool {
	return err.Code == 200
}

var (
	OK            = Errno{Code: 200, Msg: "OK"}
	ServerErr     = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Param Error"}
	AuthErr       = Errno{Code: 401, Msg: "Auth Error"}
	PermissionErr = Errno{Code: 403, Msg: "Permission Error"}

	DatabaseErr = Errno{Code: 10000, Msg: "Database Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	UserNotFoundErr   = Errno{Code: 11001, Msg: "User Not Found"}
	InvalidCaptchaErr = Errno{Code: 11002, Msg: "滑块校验失败，请重试"}

	// 位置相关错误 (12000-12999)
	LocationNotFoundErr   = Errno{Code: 12001, Msg: "Location Not Found"}
	LocationExpiredErr    = Errno{Code: 12002, Msg: "Location Data Expired"}
	InvalidCoordinatesErr = Errno{Code: 12003, Msg: "Invalid Coordinates"}

	// 好友相关错误 (13000-13999)
	FriendNotFoundErr      = Errno{Code: 13001, Msg: "Friend Not Found"}
	FriendRequestExistsErr = Errno{Code: 13002, Msg: "Friend Request Already Exists"}
	AlreadyFriendsErr      = Errno{Code: 13003, Msg: "Already Friends"}
	CannotAddSelfErr       = Errno{Code: 13004, Msg: "Cannot Add Self as Friend"}

	// 设备相关错误 (14000-14999)
	DeviceNotFoundErr     = Errno{Code: 14001, Msg: "Device Not Found"}
	DeviceAlreadyBoundErr = Errno{Code: 14002, Msg: "Device Already Bound"}
	DeviceNotBoundErr     = Errno{Code: 14003, Msg: "Device Not Bound to User"}

	// WebSocket 相关错误 (15000-15999)
	WSAuthFailedErr      = Errno{Code: 15001, Msg: "WebSocket Authentication Failed"}
	WSConnectionClosedErr = Errno{Code: 15002, Msg: "WebSocket Connection Closed"}

	// 地理围栏相关错误 (16000-16999)
	GeofenceNotFoundErr = Errno{Code: 16001, Msg: "Geofence Not Found"}
	InvalidGeofenceErr  = Errno{Code: 16002, Msg: "Invalid Geofence Parameters"}
)

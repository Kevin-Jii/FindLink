package customer

import (
	"time"

	"github.com/gin-gonic/gin"

	"app/api"
	"app/common"
	"app/service/dto"
	"app/service/location"
)

// LocationCtrl 位置控制器 - 嵌入到 Ctrl 中
type LocationCtrl struct {
	Location *location.LocationService
}

// NewLocationCtrl 创建位置控制器
func NewLocationCtrl(loc *location.LocationService) *LocationCtrl {
	return &LocationCtrl{Location: loc}
}

// @Summary 上报位置
// @Description 上报用户当前位置
// @Tags location
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.LocationReportReq true "位置信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/location/report [post]
func (c *Ctrl) Report(ctx *gin.Context) {
	req := &dto.LocationReportReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Location.ReportLocation(ctx.Request.Context(), userID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 批量上报位置
// @Description 批量上报用户位置历史
// @Tags location
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.BatchLocationReportReq true "位置列表"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/location/batch [post]
func (c *Ctrl) BatchReport(ctx *gin.Context) {
	req := &dto.BatchLocationReportReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Location.BatchReportLocation(ctx.Request.Context(), userID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, gin.H{"count": len(req.Locations)}, common.OK)
}

// @Summary 获取用户位置
// @Description 获取指定用户的位置
// @Tags location
// @Produce json
// @Param Authorization header string true "Token"
// @Param user_id path int true "用户ID"
// @Success 200 {object} api.Resp{data=dto.LocationResp}
// @Router /api/app/customer/v1/location/user/{user_id} [get]
func (c *Ctrl) GetUserLocation(ctx *gin.Context) {
	userID := getUserID(ctx)
	targetUserID := parseInt64(ctx.Param("user_id"))

	loc, err := c.Location.GetUserLocation(ctx.Request.Context(), targetUserID, userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, loc, common.OK)
}

// @Summary 获取设备位置
// @Description 获取指定设备的位置
// @Tags location
// @Produce json
// @Param Authorization header string true "Token"
// @Param device_id path string true "设备ID"
// @Success 200 {object} api.Resp{data=dto.LocationResp}
// @Router /api/app/customer/v1/location/device/{device_id} [get]
func (c *Ctrl) GetDeviceLocation(ctx *gin.Context) {
	userID := getUserID(ctx)
	deviceID := ctx.Param("device_id")

	loc, err := c.Location.GetDeviceLocation(ctx.Request.Context(), deviceID, userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, loc, common.OK)
}

// @Summary 获取位置历史
// @Description 获取用户位置历史记录
// @Tags location
// @Produce json
// @Param Authorization header string true "Token"
// @Param user_id query int true "用户ID"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Param limit query int false "限制数量"
// @Success 200 {object} api.Resp{data=[]dto.LocationResp}
// @Router /api/app/customer/v1/location/history [get]
func (c *Ctrl) GetLocationHistory(ctx *gin.Context) {
	userID := parseInt64(ctx.Query("user_id"))
	startTime := parseTime(ctx.Query("start_time"))
	endTime := parseTime(ctx.Query("end_time"))
	limit := parseInt(ctx.Query("limit"))

	req := &dto.LocationHistoryReq{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     limit,
	}

	locs, err := c.Location.GetLocationHistory(ctx.Request.Context(), req)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, locs, common.OK)
}

// @Summary 获取附近好友
// @Description 获取附近的好友列表
// @Tags location
// @Produce json
// @Param Authorization header string true "Token"
// @Param radius query number false "搜索半径（米）"
// @Success 200 {object} api.Resp{data=[]dto.NearbyFriendResp}
// @Router /api/app/customer/v1/location/nearby [get]
func (c *Ctrl) GetNearbyFriends(ctx *gin.Context) {
	userID := getUserID(ctx)
	radius := parseFloat(ctx.Query("radius"))
	if radius <= 0 {
		radius = 1000 // 默认1公里
	}

	friends, err := c.Location.GetNearbyFriends(ctx.Request.Context(), userID, radius)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, friends, common.OK)
}

// Helper functions
func getUserID(ctx *gin.Context) int64 {
	// 从JWT token中获取用户ID
	if v, exists := ctx.Get("user_id"); exists {
		return v.(int64)
	}
	return 0
}

func parseInt64(s string) int64 {
	var n int64
	for _, c := range s {
		if c < '0' || c > '9' {
			if c != '-' || n != 0 {
				break
			}
		}
		if c == '-' {
			continue
		}
		n = n*10 + int64(c-'0')
	}
	return n
}

func parseInt(s string) int {
	return int(parseInt64(s))
}

func parseFloat(s string) float64 {
	var n, frac float64
	hasFrac := false
	hasDigit := false
	for _, c := range s {
		if c == '.' {
			hasFrac = true
			continue
		}
		if c < '0' || c > '9' {
			continue
		}
		hasDigit = true
		if !hasFrac {
			n = n*10 + float64(c-'0')
		} else {
			frac = frac*10 + float64(c-'0')
		}
	}
	for frac >= 1 {
		frac /= 10
	}
	if !hasDigit {
		return 0
	}
	return n + frac
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z07:00", s)
	return t
}

package customer

import (
	"github.com/gin-gonic/gin"

	"app/api"
	"app/common"
	"app/service/dto"
	"app/service/device"
)

// DeviceCtrl 设备控制器 - 嵌入到 Ctrl 中
type DeviceCtrl struct {
	Device *device.DeviceService
}

// NewDeviceCtrl 创建设备控制器
func NewDeviceCtrl(device *device.DeviceService) *DeviceCtrl {
	return &DeviceCtrl{Device: device}
}

// @Summary 绑定设备
// @Description 绑定追踪设备
// @Tags device
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.DeviceBindReq true "设备信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/device/bind [post]
func (c *Ctrl) BindDevice(ctx *gin.Context) {
	req := &dto.DeviceBindReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Device.BindDevice(ctx.Request.Context(), userID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 解绑设备
// @Description 解绑指定设备
// @Tags device
// @Produce json
// @Param Authorization header string true "Token"
// @Param device_id path string true "设备ID"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/device/{device_id} [delete]
func (c *Ctrl) UnbindDevice(ctx *gin.Context) {
	userID := getUserID(ctx)
	deviceID := ctx.Param("device_id")

	if err := c.Device.UnbindDevice(ctx.Request.Context(), userID, deviceID); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 获取设备列表
// @Description 获取用户绑定的设备列表
// @Tags device
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} api.Resp{data=[]dto.DeviceResp}
// @Router /api/app/customer/v1/device/list [get]
func (c *Ctrl) GetDeviceList(ctx *gin.Context) {
	userID := getUserID(ctx)

	devices, err := c.Device.GetDeviceList(ctx.Request.Context(), userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, devices, common.OK)
}

// @Summary 更新设备设置
// @Description 更新设备名称和通知设置
// @Tags device
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param device_id path string true "设备ID"
// @Param req body dto.DeviceSettingsReq true "设置信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/device/{device_id}/settings [put]
func (c *Ctrl) UpdateDeviceSettings(ctx *gin.Context) {
	userID := getUserID(ctx)
	deviceID := ctx.Param("device_id")

	req := &dto.DeviceSettingsReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	if err := c.Device.UpdateDeviceSettings(ctx.Request.Context(), userID, deviceID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 更新设备状态
// @Description 更新设备电池和连接状态
// @Tags device
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param device_id path string true "设备ID"
// @Param req body dto.DeviceStatusReq true "状态信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/device/{device_id}/status [put]
func (c *Ctrl) UpdateDeviceStatus(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	req := &dto.DeviceStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	if err := c.Device.UpdateDeviceStatus(ctx.Request.Context(), deviceID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

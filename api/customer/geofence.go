package customer

import (
	"github.com/gin-gonic/gin"

	"app/api"
	"app/common"
	"app/service/dto"
	"app/service/geofence"
)

// GeofenceCtrl 地理围栏控制器 - 嵌入到 Ctrl 中
type GeofenceCtrl struct {
	Geofence *geofence.GeofenceService
}

// NewGeofenceCtrl 创建地理围栏控制器
func NewGeofenceCtrl(geofence *geofence.GeofenceService) *GeofenceCtrl {
	return &GeofenceCtrl{Geofence: geofence}
}

// @Summary 创建地理围栏
// @Description 创建新的地理围栏
// @Tags geofence
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.GeofenceCreateReq true "围栏信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/geofence [post]
func (c *Ctrl) CreateGeofence(ctx *gin.Context) {
	req := &dto.GeofenceCreateReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Geofence.Create(ctx.Request.Context(), userID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 获取地理围栏列表
// @Description 获取用户的所有地理围栏
// @Tags geofence
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} api.Resp{data=[]dto.GeofenceResp}
// @Router /api/app/customer/v1/geofence/list [get]
func (c *Ctrl) GetGeofenceList(ctx *gin.Context) {
	userID := getUserID(ctx)

	geofences, err := c.Geofence.GetList(ctx.Request.Context(), userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, geofences, common.OK)
}

// @Summary 更新地理围栏
// @Description 更新指定地理围栏
// @Tags geofence
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param geofence_id path int true "围栏ID"
// @Param req body dto.GeofenceUpdateReq true "更新信息"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/geofence/{geofence_id} [put]
func (c *Ctrl) UpdateGeofence(ctx *gin.Context) {
	userID := getUserID(ctx)
	geofenceID := parseInt64(ctx.Param("geofence_id"))

	req := &dto.GeofenceUpdateReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	if err := c.Geofence.Update(ctx.Request.Context(), userID, geofenceID, req); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 删除地理围栏
// @Description 删除指定地理围栏
// @Tags geofence
// @Produce json
// @Param Authorization header string true "Token"
// @Param geofence_id path int true "围栏ID"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/geofence/{geofence_id} [delete]
func (c *Ctrl) DeleteGeofence(ctx *gin.Context) {
	userID := getUserID(ctx)
	geofenceID := parseInt64(ctx.Param("geofence_id"))

	if err := c.Geofence.Delete(ctx.Request.Context(), userID, geofenceID); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"app/adaptor"
	appRedis "app/adaptor/redis"
	"app/api/admin"
	"app/api/customer"
	"app/common"
	"app/config"
	"app/utils/logger"
)

type IRouter interface {
	Register(engine *gin.Engine)
	SpanFilter(r *gin.Context) bool
	AccessRecordFilter(r *gin.Context) bool
}

type Router struct {
	FullPPROF bool
	rootPath  string
	conf      *config.Config
	checkFunc func() error
	admin     *admin.Ctrl
	customer  *customer.Ctrl
	verify    appRedis.IVerify
}

func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/app",
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),
		customer:  customer.NewCtrl(adaptor),
		verify:    appRedis.NewVerify(adaptor.GetRedis()),
	}
}

func (r *Router) checkServer() func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (r *Router) Register(app *gin.Engine) {
	if r.conf.Server.EnablePprof {
		SetupPprof(app, "/debug/pprof")
	}
	app.Any("/ping", r.checkServer())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	root := app.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(ctx *gin.Context) bool {
	fullPath := ctx.Request.URL.Path
	logger.Debug("SpanFilter", logger.String("fullPath", fullPath))
	_, ok := AdminAuthWhiteList[fullPath]
	return !ok
}

func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	r.customerRoute(root)
	r.adminRoute(root)
}

// validateUserToken 验证C端用户token
func (r *Router) validateUserToken(ctx context.Context, token string) (*common.User, error) {
	userID, err := r.verify.GetToken(ctx, "user:"+token)
	if err != nil {
		return nil, err
	}
	return &common.User{UserID: userID}, nil
}

// validateAdminToken 验证B端管理员token
func (r *Router) validateAdminToken(ctx context.Context, token string) (*common.AdminUser, error) {
	userID, err := r.verify.GetToken(ctx, "admin:"+token)
	if err != nil {
		return nil, err
	}
	return &common.AdminUser{UserID: userID}, nil
}

func (r *Router) customerRoute(root *gin.RouterGroup) {
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, r.validateUserToken))

	// C端用户 - 无需鉴权
	cstRoot.POST("/v1/user/login", r.customer.Login)
	cstRoot.POST("/v1/user/register", r.customer.Register)

	// C端用户 - 需要鉴权
	cstRoot.GET("/v1/user/info", r.customer.GetUserInfo)
	cstRoot.POST("/v1/user/logout", r.customer.Logout)

	// 位置相关
	locationGroup := cstRoot.Group("/v1/location")
	{
		locationGroup.POST("/report", r.customer.Report)
		locationGroup.POST("/batch", r.customer.BatchReport)
		locationGroup.GET("/user/:user_id", r.customer.GetUserLocation)
		locationGroup.GET("/device/:device_id", r.customer.GetDeviceLocation)
		locationGroup.GET("/history", r.customer.GetLocationHistory)
		locationGroup.GET("/nearby", r.customer.GetNearbyFriends)
	}

	// 好友相关
	friendGroup := cstRoot.Group("/v1/friend")
	{
		friendGroup.POST("/request", r.customer.SendFriendRequest)
		friendGroup.GET("/requests", r.customer.GetFriendRequests)
		friendGroup.POST("/accept", r.customer.AcceptFriendRequest)
		friendGroup.POST("/reject", r.customer.RejectFriendRequest)
		friendGroup.GET("/list", r.customer.GetFriendList)
		friendGroup.DELETE("/:friend_id", r.customer.RemoveFriend)
		friendGroup.GET("/search", r.customer.SearchUsers)
	}

	// 设备相关
	deviceGroup := cstRoot.Group("/v1/device")
	{
		deviceGroup.POST("/bind", r.customer.BindDevice)
		deviceGroup.GET("/list", r.customer.GetDeviceList)
		deviceGroup.PUT("/:device_id/settings", r.customer.UpdateDeviceSettings)
		deviceGroup.PUT("/:device_id/status", r.customer.UpdateDeviceStatus)
		deviceGroup.DELETE("/:device_id", r.customer.UnbindDevice)
	}

	// 地理围栏相关
	geofenceGroup := cstRoot.Group("/v1/geofence")
	{
		geofenceGroup.POST("", r.customer.CreateGeofence)
		geofenceGroup.GET("/list", r.customer.GetGeofenceList)
		geofenceGroup.PUT("/:geofence_id", r.customer.UpdateGeofence)
		geofenceGroup.DELETE("/:geofence_id", r.customer.DeleteGeofence)
	}

	// WebSocket
	cstRoot.GET("/v1/ws", r.customer.WebSocketConnect)
}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, r.validateAdminToken))

	// 登录无鉴权：添加白名单
	adminRoot.GET("/v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	adminRoot.POST("/v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)
	adminRoot.POST("/v1/user/login", r.admin.Login)

	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
}

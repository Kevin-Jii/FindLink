package router

import (
	"app/adaptor"
	"app/api/admin"
	"app/api/customer"
	"app/common"
	"app/config"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
}

func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/app",
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),
		customer:  customer.NewCtrl(adaptor),
	}
}

func (r *Router) checkServer() func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
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

	// Swagger 文档
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	root := app.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(ctx *gin.Context) bool {
	path := strings.Replace(ctx.Request.URL.Path, r.rootPath, "", 1)
	_, ok := AdminAuthWhiteList[path]
	if ok {
		return false
	}
	return true
}

func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	r.customerRoute(root)
	r.adminRoute(root)
}

func (r *Router) customerRoute(root *gin.RouterGroup) {
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.User, error) {
		return &common.User{}, nil
	}))

	// C端用户 - 无需鉴权
	cstRoot.POST("/v1/user/login", r.customer.Login)
	cstRoot.POST("/v1/user/register", r.customer.Register)

	// C端用户 - 需要鉴权
	cstRoot.GET("/v1/user/info", r.customer.GetUserInfo)
}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.AdminUser, error) {
		return &common.AdminUser{
			UserID: 1,
			Name:   "admin",
		}, nil
	}))
	// 登录无鉴权：添加白名单
	adminRoot.GET("/v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	adminRoot.POST("/v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)
	adminRoot.POST("/v1/user/login", r.admin.Login)

	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
}

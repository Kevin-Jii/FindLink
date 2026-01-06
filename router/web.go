package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"app/middleware"
	"app/utils/logger"
)

type App struct {
	server          *gin.Engine
	addr            string
	shutdownTimeout time.Duration
}

func NewApp(port int, shutdownTimeout int, router IRouter) *App {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 中间件顺序很重要
	engine.Use(middleware.TraceMiddleware())                   // 1. 请求追踪
	engine.Use(middleware.RecoveryMiddleware())                // 2. Panic恢复
	engine.Use(corsMiddleware())                               // 3. CORS
	engine.Use(AccessLogMiddleware(router.AccessRecordFilter)) // 4. 访问日志

	// 注册业务路由
	router.Register(engine)

	timeout := time.Duration(shutdownTimeout) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &App{
		server:          engine,
		addr:            ":" + strconv.Itoa(port),
		shutdownTimeout: timeout,
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Trace-ID, token")
		c.Header("Access-Control-Expose-Headers", "X-Trace-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Run 启动服务并支持优雅关闭
func (app *App) Run() {
	srv := &http.Server{
		Addr:         app.addr,
		Handler:      app.server,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 异步启动服务
	go func() {
		logger.Info(fmt.Sprintf("server started, endpoint: http://localhost%s", app.addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Warn("received shutdown signal", zap.String("signal", sig.String()))
	logger.Info(fmt.Sprintf("shutting down server with %v timeout...", app.shutdownTimeout))

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), app.shutdownTimeout)
	defer cancel()

	// 优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("server exited gracefully")
	}
}

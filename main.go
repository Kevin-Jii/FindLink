package main

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"app/adaptor"
	"app/adaptor/repo/model"
	"app/config"
	_ "app/docs"
	"app/router"
	"app/utils/logger"
)

// @title           App API
// @version         1.0
// @description     基础后端服务API文档

// @host      localhost:8080
// @BasePath  /api/app

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token

func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.Server.LogLevel)

	dbClient, err := initMysql(&conf.Mysql)
	handleErr(err)
	logger.Debug("mysql connect success")

	rdsClient, err := initRedis(&conf.Redis)
	handleErr(err)
	logger.Debug("client connect success")

	startServer(conf, dbClient, rdsClient).Run()
}

func startServer(conf *config.Config, db *gorm.DB, redis *redis.Client) *router.App {
	return router.NewApp(conf.Server.HttpPort,
		router.NewRouter(
			conf,
			adaptor.NewAdaptor(conf, db, redis),
			func() error {
				err := func() error {
					pingDb, err := db.DB()
					handleErr(err)
					return pingDb.Ping()
				}()
				if err != nil {
					return errors.New("mysql connect failed")
				}
				return redis.Ping().Err()
			},
		),
	)
}

func initRedis(conf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.PWD,
		DB:           conf.DBIndex,
		MinIdleConns: conf.MaxIdle,
		PoolSize:     conf.MaxOpen,
	})
	if r, _ := client.Ping().Result(); r != "PONG" {
		return nil, errors.New("redis connect failed")
	}
	return client, nil
}

func initMysql(conf *config.Mysql) (*gorm.DB, error) {
	conf.MaxIdle = lo.Max([]int{conf.MaxIdle + 1, 5})
	conf.MaxOpen = lo.Max([]int{conf.MaxOpen + 1, 10})
	dsn := conf.GetDsn()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetMaxOpenConns(conf.MaxOpen)

	// 自动迁移表结构
	if err = autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.AdminUser{},
		&model.AdminUserRole{},
		&model.Permission{},
		&model.Role{},
		&model.RolePermission{},
	)
	if err != nil {
		return err
	}

	// 初始化超级管理员账户
	initSuperAdmin(db)
	return nil
}

func initSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.AdminUser{}).Where("name = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	// 创建超级管理员 密码: admin123
	admin := &model.AdminUser{
		Name:     "admin",
		NickName: "超级管理员",
		Password: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9", // sha256("admin123")
		Status:   1,
		Sex:      1,
	}
	db.Create(admin)
	logger.Info("初始化超级管理员账户成功: admin / admin123")
}
